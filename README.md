# MQTT to PostgreSQL Integration with Node.js

## Project Overview
This project demonstrates how to create a backend service using Node.js to process MQTT messages, store the data in a PostgreSQL database, and expose RESTful APIs for querying the stored data. The system handles JSON payloads from MQTT topics, stores the data efficiently, and provides endpoints to retrieve, filter, and analyze the data.

---

## Features
- **MQTT Integration**: Subscribes to MQTT topics to handle incoming JSON messages.
- **Data Storage**: Stores data in a PostgreSQL database with a well-defined schema.
- **RESTful APIs**: Provides endpoints for querying the latest data, historical data, and average metrics.

---

## Swagger API Documentation
This project uses Swagger to provide interactive API documentation. You can access the Swagger UI to explore and test the API endpoints.

### Accessing Swagger UI
1. Start the server by running:
   ```bash
   npm start
   ```

---

## Requirements
- **Node.js** (v14 or later)
- **PostgreSQL**
- **MQTT Broker** (e.g., Mosquitto, HiveMQ)

### Dependencies
- [pg](https://www.npmjs.com/package/pg): PostgreSQL client for Node.js.
- [mqtt](https://www.npmjs.com/package/mqtt): MQTT client library.
- [dotenv](https://www.npmjs.com/package/dotenv): For managing environment variables.
- [express](https://www.npmjs.com/package/express): To create RESTful APIs.

---

## Installation

### Step 1: Clone the Repository
```bash
https://github.com/RafiIkhwan/mqtt-backend-service.git
cd mqtt-backend-service
```

### Step 2: Install Dependencies
```bash
npm install
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

## Getting Started

These instructions will get you a copy of the project up and running on your local machine for development and testing purposes. See deployment for notes on how to deploy the project on a live system.

## MakeFile

Run build make command with tests
```bash
make all
```

Build the application
```bash
make build
```

Run the application
```bash
make run
```
Create DB container
```bash
make docker-run
```

Shutdown DB Container
```bash
make docker-down
```

DB Integrations Test:
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

Clean up binary from the last build:
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
├── app/
│   ├── config/
│   │   ├── db.js           # PostgreSQL connection setup
│   │   └── mqtt_broker.js  # MQTT broker configuration
│   ├── routes/
│   │   └── devices_data.js # API routes for device data
├── src/
│   └── subscribe.js        # Handles MQTT subscriptions and message processing
├── .env                    # Environment variable configuration
├── .env.example            # Example environment variables
├── .gitignore              # Git ignore file
├── index.js                # Entry point of the application
├── package.json            # Project dependencies and scripts
├── package-lock.json       # Dependency lockfile
└── README.md               # Project documentation
```

---

## Usage

### Start the Server
```bash
npm start
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

#### 2. Get History Data
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