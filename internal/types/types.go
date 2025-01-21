package types

import "time"

type DeviceData struct {
    DeviceID    string    `json:"device_id"`
    Humidity    float64   `json:"humidity"`
    Temperature float64   `json:"temperature"`
    Timestamp   time.Time `json:"timestamp"`
}

type AverageData struct {
		AverageHumidity    	float64    `json:"average_humidity"`
    AverageTemperature  float64   `json:"average_temperature"`
}