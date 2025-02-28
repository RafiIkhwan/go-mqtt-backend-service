package database

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	_ "github.com/joho/godotenv/autoload"
	
    "mqtt-backend-service/internal/types"
)

// Service represents a service that interacts with a database.
type Service interface {
	// Health returns a map of health status information.
	// The keys and values in the map are service-specific.
	Health() map[string]string

	GetLatestData() ([]types.DeviceData, error)
	GetHistoryData(deviceID string, startDate, endDate time.Time) ([]types.DeviceData, error)
	GetAverageData(deviceID string, startDate, endDate time.Time) (types.AverageData, error)

	InsertDeviceData(data types.DeviceData) error
	// Close terminates the database connection.
	// It returns an error if the connection cannot be closed.
	Close() error
}

type service struct {
	db *sql.DB
}

var (
	database   = os.Getenv("BLUEPRINT_DB_DATABASE")
	password   = os.Getenv("BLUEPRINT_DB_PASSWORD")
	username   = os.Getenv("BLUEPRINT_DB_USERNAME")
	port       = os.Getenv("BLUEPRINT_DB_PORT")
	host       = os.Getenv("BLUEPRINT_DB_HOST")
	schema     = os.Getenv("BLUEPRINT_DB_SCHEMA")
	dbInstance *service
)

func New() Service {
	// Reuse Connection
	if dbInstance != nil {
		return dbInstance
	}
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable&search_path=%s", username, password, host, port, database, schema)
	db, err := sql.Open("pgx", connStr)
	if err != nil {
		log.Fatal(err)
	}
	dbInstance = &service{
		db: db,
	}
	return dbInstance
}

// Health checks the health of the database connection by pinging the database.
// It returns a map with keys indicating various health statistics.
func (s *service) Health() map[string]string {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	stats := make(map[string]string)

	// Ping the database
	err := s.db.PingContext(ctx)
	if err != nil {
		stats["status"] = "down"
		stats["error"] = fmt.Sprintf("db down: %v", err)
		log.Fatalf("db down: %v", err) // Log the error and terminate the program
		return stats
	}

	// Database is up, add more statistics
	stats["status"] = "up"
	stats["message"] = "It's healthy"

	// Get database stats (like open connections, in use, idle, etc.)
	dbStats := s.db.Stats()
	stats["open_connections"] = strconv.Itoa(dbStats.OpenConnections)
	stats["in_use"] = strconv.Itoa(dbStats.InUse)
	stats["idle"] = strconv.Itoa(dbStats.Idle)
	stats["wait_count"] = strconv.FormatInt(dbStats.WaitCount, 10)
	stats["wait_duration"] = dbStats.WaitDuration.String()
	stats["max_idle_closed"] = strconv.FormatInt(dbStats.MaxIdleClosed, 10)
	stats["max_lifetime_closed"] = strconv.FormatInt(dbStats.MaxLifetimeClosed, 10)

	// Evaluate stats to provide a health message
	if dbStats.OpenConnections > 40 { // Assuming 50 is the max for this example
		stats["message"] = "The database is experiencing heavy load."
	}

	if dbStats.WaitCount > 1000 {
		stats["message"] = "The database has a high number of wait events, indicating potential bottlenecks."
	}

	if dbStats.MaxIdleClosed > int64(dbStats.OpenConnections)/2 {
		stats["message"] = "Many idle connections are being closed, consider revising the connection pool settings."
	}

	if dbStats.MaxLifetimeClosed > int64(dbStats.OpenConnections)/2 {
		stats["message"] = "Many connections are being closed due to max lifetime, consider increasing max lifetime or revising the connection usage pattern."
	}

	return stats
}


func (s *service) GetLatestData() ([]types.DeviceData, error) {
	// MIGRATION IF TABLE NOT EXISTS
	_, err := s.db.Exec("CREATE TABLE IF NOT EXISTS device_data (device_id TEXT, humidity FLOAT, temperature FLOAT, timestamp TIMESTAMP)")
	if err != nil {
		log.Fatal(err)
	}

	rows, err := s.db.Query(`
			SELECT DISTINCT ON (device_id) device_id, humidity, temperature, timestamp
			FROM device_data
			ORDER BY device_id, timestamp DESC
	`)
	if err != nil {
			return nil, err
	}
	defer rows.Close()

	var data []types.DeviceData
	for rows.Next() {
			var d types.DeviceData
			if err := rows.Scan(&d.DeviceID, &d.Humidity, &d.Temperature, &d.Timestamp); err != nil {
					return nil, err
			}
			data = append(data, d)
	}
	return data, nil
}


func (s *service) GetHistoryData(deviceID string, startDate, endDate time.Time) ([]types.DeviceData, error) {
	rows, err := s.db.Query(
			`SELECT device_id, humidity, temperature, timestamp FROM device_data WHERE device_id = $1 AND timestamp BETWEEN $2 AND $3`,
			deviceID, startDate, endDate,
	)
	if err != nil {
			return nil, err
	}
	defer rows.Close()

	var data []types.DeviceData
	for rows.Next() {
			var d types.DeviceData
			if err := rows.Scan(&d.DeviceID, &d.Humidity, &d.Temperature, &d.Timestamp); err != nil {
					return nil, err
			}
			data = append(data, d)
	}
	return data, nil
}

func (s *service) GetAverageData(deviceID string, startDate, endDate time.Time) (types.AverageData, error) {
	row := s.db.QueryRow(
			`SELECT 
			ROUND(CAST(AVG(humidity) AS numeric), 2) AS average_humidity, 
			ROUND(CAST(AVG(temperature) AS numeric), 2) AS average_temperature
			FROM device_data
			WHERE device_id = $1 AND timestamp BETWEEN $2 AND $3`,
			deviceID, startDate, endDate,
	)

	var avgData types.AverageData
	if err := row.Scan(&avgData.AverageHumidity, &avgData.AverageTemperature); err != nil {
			return avgData, err
	}
	return avgData, nil
}


func (s *service) InsertDeviceData(data types.DeviceData) error {
	_, err := s.db.Exec(
			`INSERT INTO device_data (device_id, humidity, temperature, timestamp) VALUES ($1, $2, $3, $4)`,
			data.DeviceID, data.Humidity, data.Temperature, data.Timestamp,
	)
	return err
}

// Close closes the database connection.
// It logs a message indicating the disconnection from the specific database.
// If the connection is successfully closed, it returns nil.
// If an error occurs while closing the connection, it returns the error.
func (s *service) Close() error {
	log.Printf("Disconnected from database: %s", database)
	return s.db.Close()
}
