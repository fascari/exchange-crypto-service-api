.DEFAULT_GOAL := help

help:
	@echo "API SERVICE MAKEFILE"
	@echo ""
	@echo "🐳 DOCKER COMMANDS:"
	@echo " run-docker    Run application via Docker containers"
	@echo " run-db        Run only PostgreSQL database"
	@echo " build         Build containers with no cache"
	@echo " fs            Fresh start containers"
	@echo " down          Stop and remove containers"
	@echo " restart       Restart API container"
	@echo " logs          Show API container logs"
	@echo " status        Show containers status"
	@echo ""
	@echo "🏃 DEVELOPMENT COMMANDS:"
	@echo " run-local     Run application locally (Go binary)"
	@echo " migrate       Run database migration"
	@echo ""
	@echo "🧪 TESTING & QUALITY:"
	@echo " test          Run unit tests"
	@echo " lint-code     Run golangci-lint"
	@echo ""

run-db:
	@echo "Starting PostgreSQL database..."
	docker compose up -d postgres
	@echo "✅ PostgreSQL is ready at localhost:5435"

run-local:
	@echo "Starting dependencies..."
	docker compose up -d postgres otel-collector jaeger
	@echo "Waiting for database to be ready..."
	@sleep 5
	@$(MAKE) migrate
	@echo "Starting the application locally..."
	go run cmd/api/*.go

run-docker:
	@echo "Starting all services via Docker..."
	docker compose up -d
	@echo "API is running at http://localhost:8080"
	@echo "Jaeger UI is available at http://localhost:16686"

build:
	docker compose build --no-cache --force-rm --pull

fs:
	docker compose down
	docker compose up -d
	@echo "Fresh start completed. API is running at http://localhost:8080"

restart:
	@echo "Restarting API container..."
	docker compose restart exchange-crypto-api
	@echo "API container restarted"

logs:
	@echo "Showing API logs..."
	docker compose logs -f exchange-crypto-api

status:
	@echo "Containers status:"
	docker compose ps

test:
	@echo "🧪 Running tests..."
	@go test -v -short $(shell go list ./internal/... ./pkg/... | grep -v -E 'test|mocks') -shuffle=on -failfast -coverprofile=coverage.out
	@echo "📊 Generating coverage report..."
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

migrate:
	@docker compose run --rm liquibase liquibase --logLevel=warning --defaultsFile=/liquibase/changelog/local.properties update
	@echo "✅ Migration completed"

lint-code:
	@echo "Running golangci-lint..."
	@type "golangci-lint" > /dev/null 2>&1 || (echo 'Please install golangci-lint' && exit 1)
	@golangci-lint run --config=.golangci.yml --verbose ./...