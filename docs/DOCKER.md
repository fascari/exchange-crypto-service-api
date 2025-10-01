# 🐳 Docker Guide

This application is fully containerized and can be run entirely via Docker containers. All services including the API, database, observability tools, and migrations are orchestrated using Docker Compose.

## Prerequisites
- Docker (>= 20.10)
- Docker Compose (>= 2.0)

## Running with Docker

### Start All Services
```bash
# Start the complete application stack
make run-docker
```
This command starts:
- **PostgreSQL Database** (port 5435)
- **Exchange Crypto API** (port 8080)
- **Liquibase** (automatic database migration)
- **OpenTelemetry Collector** (ports 4317, 4318, 8889)
- **Jaeger UI** (port 16686 for tracing visualization)

### Docker Management Commands
```bash
# Build containers from scratch (no cache)
make build

# Fresh start - stop everything and restart
make fs

# Stop all containers
make down

# Restart only the API container
make restart

# View API logs in real-time
make logs

# Check containers status
make status
```

## Service URLs
After running `make run-docker`, the following services will be available:

| Service | URL | Description |
|---------|-----|-------------|
| **API** | http://localhost:8080 | Main Exchange Crypto API |
| **Health Check** | http://localhost:8080/health | API health status |
| **Jaeger UI** | http://localhost:16686 | Distributed tracing dashboard |
| **PostgreSQL** | localhost:5435 | Database (user: owner, pass: owner123) |

## Environment Variables

The API container uses the following environment variables (defined in `docker-compose.yml`):

```yaml
environment:
  - DB_URL=postgres
  - DB_PORT=5432
  - DB_APP_USER=owner
  - DB_APP_PASSWORD=owner123
  - DB_NAME=exchange_crypto_local
  - DB_SSL_MODE=disable
  - DB_SCHEMA=exchange_crypto
  - DB_MAX_IDLE_CONNECTIONS=10
  - DB_MAX_OPEN_CONNECTIONS=100
  - DB_MAX_LIFETIME_CONNECTIONS=60
```