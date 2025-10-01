# Development Guide

## Testing

### Unit Tests
Run all unit tests with coverage:
```bash
make test
```

### Code Linting with golangci-lint

The project uses [golangci-lint](https://golangci-lint.run/) for code linting and static analysis to maintain code quality and consistency.

#### Installation

**Using Go install:**
```bash
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
```

**Using Homebrew (macOS):**
```bash
brew install golangci-lint
```

**Using Binary (Linux/Windows):**
```bash
curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin
```

#### Configuration

The project includes a `.golangci.yml` configuration file that defines:
- **Enabled linters**: Code quality, style, and security checks
- **Disabled rules**: Project-specific exceptions
- **File exclusions**: Generated files, vendor directories
- **Severity levels**: Error vs warning classifications

#### Running Linter

To run the linter on the entire project:
```bash
make lint-code
```

To run golangci-lint directly:
```bash
# Run with default configuration
golangci-lint run

# Run with custom config file
golangci-lint run --config=.golangci.yml

# Run on specific directories
golangci-lint run ./internal/... ./pkg/...

# Run with verbose output
golangci-lint run --verbose
```

#### Common Linting Issues

- **Unused variables/imports**: Remove or use underscore `_` for intentionally unused
- **Naming conventions**: Use camelCase for private, PascalCase for public
- **Error handling**: Always check and handle errors appropriately
- **Code complexity**: Keep functions simple and focused
- **Documentation**: Add comments for exported functions and types

### Mockery - Mock Generation

The project uses [Mockery](https://vektra.github.io/mockery/) to automatically generate mocks for interfaces used in testing.

#### Installation
```bash
go install github.com/vektra/mockery/v2@latest
```

#### Configuration
The project includes a `.mockery.yaml` configuration file:
```yaml
with-expecter: true
case: snake
disable-version-string: true
```

#### Generating Mocks
To generate mocks for all interfaces:
```bash
mockery --all
```

To generate mocks for a specific interface:
```bash
mockery --name=Repository --dir=./internal/app/account/usecase/createaccount
```

#### Using Generated Mocks
The mocks are automatically generated in `mocks/` directories within each package. Example usage:

```go
// Create mock instance
mockRepo := mocks.NewRepository(t)

// Set expectations using EXPECT() syntax
mockRepo.EXPECT().Create(mock.Anything, mock.AnythingOfType("domain.Account")).Return(expectedAccount, nil)

// Use mock in tests
useCase := createaccount.New(mockRepo, mockUserRepo, mockExchangeRepo)
result, err := useCase.Execute(context.Background(), inputAccount)
```

#### Mock Location
- **Handler tests**: Use mocks from `usecase/<module>/mocks/`
- **UseCase tests**: Mock repository interfaces directly
- **Repository tests**: Use real database (SQLite in-memory for tests)
