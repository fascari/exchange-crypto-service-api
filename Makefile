.DEFAULT_GOAL := help

help:
	@echo "API SERVICE MAKEFILE"
	@echo ""
	@echo "COMMANDS:"
	@echo " run-local     Run application locally"
	@echo " fs            Fresh start containers"
	@echo " build         Build containers with no cache"
	@echo " down          Stop and remove containers"
	@echo " test          Run unit tests"
	@echo " lint-code     Run golangci-lint"
	@echo " migrate       Run database migration"
	@echo ""

run-local:
	@echo "Starting the application"
	docker-compose up -d
	go run cmd/api/*.go

build:
	docker-compose build --no-cache --force-rm --pull

fs:
	docker-compose down
	docker-compose up -d
	@$(MAKE) migration

down:
	docker-compose down

test:
	@echo "🧪 Running tests..."
	@go test -v -short $(shell go list ./internal/... ./pkg/... | grep -v -E 'test|mocks') -shuffle=on -failfast -coverprofile=coverage.out
	@echo "📊 Generating coverage report..."
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

migrate:
	docker-compose run --rm liquibase liquibase --logLevel=info --defaultsFile=/liquibase/changelog/local.properties update

lint-code:
	@echo "Running golangci-lint..."
	@type "golangci-lint" > /dev/null 2>&1 || (echo 'Please install golangci-lint' && exit 1)
	@golangci-lint run --config=.golangci.yml --verbose ./...