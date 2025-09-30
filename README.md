# Exchange Crypto Service API

A API for cryptocurrency exchange operations, built in Go, following Clean Architecture principles with integrated observability.

## Tech Stack
- **Go** (>=1.25)
- **PostgreSQL** (data persistence)
- **Liquibase** (database migration)
- **Docker & Docker Compose** (containerization and orchestration)
- **OpenTelemetry** (tracing and observability)
- **Makefile** (command automation)

## Project Structure
- `cmd/api/main.go`: application entry point
- `internal/app`: business modules (account, exchange, transaction, user)
- `internal/config`: configurations
- `internal/database`: database connection
- `liquibase/`: database scripts and migrations
- `docker-compose.yml`: services orchestration
- `Makefile`: automated commands

## Configuration
- Environment variables can be defined in `env.yaml`

## Database
- PostgreSQL is used as the main database.
- Automatic migrations via Liquibase (`liquibase/changelog/migrations/*.sql`).

## Basic Commands

### Using Makefile
```bash
# Run application locally (starts containers + API)
make run-local

# Fresh start (rebuild containers and run migrations)
make fs

# Build containers
make build

# Run database migrations
make migrate

# Run tests
make test

# Run code linting
make lint-code

# Stop all services
make down

# Show help with available commands
make help
```

## Main Endpoints
- **POST /api/v1/users**: create user
- **GET /api/v1/users/{userId}/balance**: check balance
- **POST /api/v1/accounts**: create account
- **POST /api/v1/transactions**: create transaction (deposit/withdrawal)
- **GET /api/v1/transactions/daily?date=YYYY-MM-DD**: get daily transactions
- **GET /health**: health check

## Telemetry & Observability

The application includes comprehensive observability through OpenTelemetry integration with Jaeger for distributed tracing.

### Architecture
- **OpenTelemetry Collector**: Collects and processes telemetry data
- **Jaeger**: Stores and visualizes distributed traces
- **Automatic Instrumentation**: HTTP requests, database operations, and custom spans

### Accessing Jaeger UI

Once you start the application with `make run-local`, the Jaeger UI will be available at:

**http://localhost:16686**
