# CloakApp Receiver

A high-performance telemetry receiver built with Go (Gin) and InfluxDB.

## Stack
- **Backend:** Go (Gin Framework)
- **Database:** InfluxDB (Time-series)
- **Frontend Integration:** REST API (JSON)

## Getting Started

### 1. Prerequisites
- Go 1.25+
- InfluxDB v2 instance

### 2. Setup
1. Copy `.env.example` to `.env` (or set the environment variables).
   ```bash
   cp .env.example .env
   ```
2. Install dependencies:
   ```bash
   go mod tidy
   ```

### 3. Run with Podman/Docker Compose
```bash
podman-compose up --build
# OR
docker-compose up --build
```
This will start both the Go receiver and InfluxDB. InfluxDB will be automatically initialized with the credentials defined in `compose.yml`.

## API Reference

### POST `/telemetry`
Receives JSON telemetry data.

**Payload Structure:**
```json
{
  "host": "example.com",
  "package": "com.domain.test",
  "full_url": "https://example.com/page?id=123",
  "data": {
    "event": "page_load",
    "load_time_ms": 250
  }
}
```
*Note: `ip` and `user_agent` are automatically detected by the server and stored as tags for each entry.*

## Testing
You can use `client_example.js` in your browser's console or include it in your frontend project to test the integration.
