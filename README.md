# MQTT to PostgreSQL Integration with Go

## Project Overview
This project demonstrates how to create a backend service using Go to process MQTT messages, store the data in a PostgreSQL database, and expose RESTful APIs for querying the stored data. The system handles JSON payloads from MQTT topics, stores the data efficiently, and provides endpoints to retrieve, filter, and analyze the data.

---

## Features
- **MQTT Integration**: Subscribes to MQTT topics to handle incoming JSON messages.
- **Data Storage**: Stores data in a PostgreSQL database with a well-defined schema.
- **RESTful APIs**: Provides endpoints for querying the latest data, historical data, and average metrics.

---

## Requirements
- **Go** (v1.19 or later)
- **PostgreSQL**
- **MQTT Broker** (e.g., Mosquitto, HiveMQ)

### Dependencies
This project uses the following Go packages:
- [github.com/eclipse/paho.mqtt.golang](https://pkg.go.dev/github.com/eclipse/paho.mqtt.golang): MQTT client library.
- [github.com/lib/pq](https://pkg.go.dev/github.com/lib/pq): PostgreSQL client library.
- [github.com/gin-gonic/gin](https://pkg.go.dev/github.com/gin-gonic/gin): Web framework for building APIs.
- [github.com/joho/godotenv](https://pkg.go.dev/github.com/joho/godotenv): Loads environment variables from a `.env` file.

---

## Installation

### Step 1: Clone the Repository
```bash
https://github.com/RafiIkhwan/mqtt-backend-service.git
cd mqtt-backend-service
```

### Step 2: Install Dependencies
```bash
go mod tidy
```

### Step 3: Set Up Environment Variables
Create a `.env` file in the root directory with the following variables:
```env
# Web Configurations
PORT=8080
APP_ENV=local

# PostgreSQL Configurations
BLUEPRINT_DB_HOST=localhost
BLUEPRINT_DB_PORT=5432
BLUEPRINT_DB_DATABASE=mydatabase
BLUEPRINT_DB_USERNAME=myuser  # Leave empty if not required
BLUEPRINT_DB_PASSWORD=mypassword  # Leave empty if not required
BLUEPRINT_DB_SCHEMA=public

# MQTT Configurations
MQTT_BROKER=mqtt://broker.emqx.io
MQTT_TOPIC=devices/data
```

### Step 4: Run the Application
Use the `make` commands provided in the `Makefile` to build and run the application.

---

## Makefile Commands

Run build and tests:
```bash
make all
```

Build the application:
```bash
make build
```

Run the application:
```bash
make run
```

Create a PostgreSQL container:
```bash
make docker-run
```

Shut down the PostgreSQL container:
```bash
make docker-down
```

Run integration tests:
```bash
make itest
```

Live reload the application:
```bash
make watch
```

Run the test suite:
```bash
make test
```

Clean up binaries from the last build:
```bash
make clean
```

---

## Database Schema
Run the following SQL command to create the required table:
```sql
CREATE TABLE device_data (
    id SERIAL PRIMARY KEY,
    device_id VARCHAR(255),
    humidity NUMERIC,
    temperature NUMERIC,
    timestamp TIMESTAMP
);
```

---

## Project Structure
```plaintext
mqtt-backend-service/
├── cmd/
│   └── api/
│       └── main.go          # Entry point for the application
├── internal/
│   ├── database/
│   │   ├── database.go      # Database connection and queries
│   │   └── database_test.go # Database-related tests
│   ├── server/
│   │   ├── routes.go        # API routes configuration
│   │   ├── routes_test.go   # API route tests
│   │   └── server.go        # Server setup and initialization
│   ├── subscribe/
│   │   └── subscribe.go     # Handles MQTT subscriptions and message processing
│   └── types/
│       └── types.go         # Shared data structures and types
├── .air.toml                # Live reload configuration
├── .env                     # Environment variable configuration
├── .env.example             # Example environment variables
├── .gitignore               # Git ignore file
├── docker-compose.yml       # Docker Compose configuration
├── dockerfile               # Dockerfile for the application
├── go.mod                   # Go module dependencies
├── go.sum                   # Dependency lockfile
├── Makefile                 # Makefile for running tasks
└── README.md                # Project documentation
```

---

## Usage

### Start the Server
```bash
make run
```

The server will start and subscribe to the MQTT topic defined in the `.env` file.

### API Endpoints

#### 1. Get Latest Data
**Endpoint**: `GET /api/data/latest`
- **Description**: Returns the most recent data entry for all devices.
- **Response Example**:
```json
[
  {
    "device_id": "device1",
    "humidity": 55.3,
    "temperature": 23.1,
    "timestamp": "2025-01-14T10:30:00Z"
  }
]
```

#### 2. Get Historical Data
**Endpoint**: `GET /api/data/history`
- **Query Parameters**:
  - `start`: Start date-time (ISO 8601).
  - `end`: End date-time (ISO 8601).
  - `device_id`: (Optional) Filter by device ID.
- **Response Example**:
```json
[
  {
    "device_id": "device1",
    "humidity": 55.3,
    "temperature": 23.1,
    "timestamp": "2025-01-14T10:30:00Z"
  }
]
```

#### 3. Get Average Data
**Endpoint**: `GET /api/data/average`
- **Query Parameters**:
  - `start`: Start date-time (ISO 8601).
  - `end`: End date-time (ISO 8601).
  - `device_id`: (Optional) Filter by device ID.
- **Response Example**:
```json
{
  "average_humidity": 57.75,
  "average_temperature": 23.95
}
```

---

## Error Handling
- **Database Connection Errors**: Ensure PostgreSQL service is running and the credentials in the `.env` file are correct.
- **MQTT Connection Errors**: Verify the MQTT broker URL and topic in the `.env` file.

---

## License
This project is licensed under the MIT License.

---

## Contributing
Pull requests are welcome. For significant changes, please open an issue to discuss what you would like to change.

---

## Contact
If you have any questions or need further assistance, feel free to reach out:
- **Email**: rafiikhwan2006@gmail.com
- **Portfolio**: [rafiikhwan.my.id](http://rafiikhwan.my.id)