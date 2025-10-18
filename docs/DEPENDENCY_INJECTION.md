# Dependency Injection with Uber/FX

## Overview

The application uses **uber/fx** for dependency injection, providing:

- **Automatic dependency resolution** - No manual wiring needed
- **Lifecycle management** - Graceful startup and shutdown
- **Modular architecture** - Clear separation of concerns
- **Easy testing** - Simple to replace dependencies with mocks

## Architecture Overview

```
┌─────────────────────────────────────────────────────────────┐
│                         main.go                             │
│                    (Application Entry)                      │
└────────────────────┬────────────────────────────────────────┘
                     │
                     │ Uber/FX Dependency Injection
                     │
        ┌────────────┴─────────────┐
        │                          │
        ▼                          ▼
┌──────────────────┐       ┌──────────────────┐
│  Infrastructure  │       │ Business Modules │
│(internal/bootstrap)      │ (cmd/api/modules)│
├──────────────────┤       ├──────────────────┤
│ • Config         │       │ • Health         │
│ • Database       │       │ • Auth           │
│ • Telemetry      │       │ • User           │
│ • Router         │       │ • Account        │
│ • Server         │       │ • Transaction    │
│ • Logger         │       │ • Exchange       │
└────────┬─────────┘       └────────┬─────────┘
         │                          │
         │     Each Module Has:     │
         │     ┌──────────────┐     │
         └────►│ Repository   │◄────┘
               │ UseCase      │
               │ Handler      │
               └──────────────┘
```

**Key Design Principles:**
- **Separation of Concerns**: Infrastructure (`internal/bootstrap`) vs Business Logic (`cmd/api/modules`)
- **Dependency Injection**: Uber/FX manages all dependencies automatically
- **Modular Architecture**: Each domain is self-contained and independent
- **Clean Architecture**: Business logic is independent of frameworks and infrastructure

## Configuration

### FX_VERBOSE Environment Variable

By default, Uber/FX logs are disabled to keep the application output clean. You can enable verbose FX logs for debugging dependency injection issues:

```bash
# Enable FX verbose logs
FX_VERBOSE=true go run cmd/api/main.go

# Or with make
FX_VERBOSE=true make run-local
```

**When to use:**
- ✅ Debugging dependency injection issues
- ✅ Understanding the startup sequence
- ✅ Troubleshooting circular dependencies
- ❌ Production environments (keep it `false`)

**In Docker (`docker-compose.yml`):**
```yaml
environment:
  - FX_VERBOSE=true  # change to true for debugging
```

## Module Structure

### Infrastructure Modules (`internal/bootstrap/`)

| Module | Responsibility |
|--------|----------------|
| `Config` | Application configurations (Database, JWT, RateLimiter) |
| `Database` | PostgreSQL connection via GORM |
| `Telemetry` | OpenTelemetry tracer with lifecycle hooks |
| `Router` | HTTP routers and middleware setup |
| `Server` | HTTP server with graceful shutdown |
| `Logger` | FX logger configuration |

### Business Modules (`cmd/api/modules/`)

| Module | Responsibility |
|--------|----------------|
| `Health` | Health check endpoint |
| `Auth` | Authentication and token generation |
| `User` | User management (create, find balance) |
| `Account` | Account management |
| `Transaction` | Transaction operations (create, find daily) |
| `Exchange` | Exchange repository |

## How It Works

### Application Bootstrap

```go
func main() {
    logger.Init()
    jwt.Initialize()

    app := fx.New(
        bootstrap.Logger(),
        // Infrastructure
        bootstrap.Config,
        bootstrap.Database,
        bootstrap.Telemetry,
        bootstrap.Router,
        bootstrap.Server,
        // Business Modules
        modules.Health,
        modules.Exchange,
        modules.User,
        modules.Auth,
        modules.Account,
        modules.Transaction,
    )

    app.Run()
}
```

### Dependency Chain Example

```
Database → Repository → UseCase → Handler
```

Uber/fx automatically resolves and injects all dependencies in the correct order.

## Adding New Components

### Creating a New Business Module

Create a new file `/cmd/api/modules/mymodule.go`:

```go
package modules

import (
    myrepo "exchange-crypto-service-api/internal/app/mymodule/repository"
    myusecase "exchange-crypto-service-api/internal/app/mymodule/usecase"
    "exchange-crypto-service-api/internal/app/mymodule/handler"
    
    "go.uber.org/fx"
    "gorm.io/gorm"
)

var MyModule = fx.Module("mymodule",
    fx.Provide(
        NewMyRepository,
        NewMyUseCase,
    ),
    fx.Invoke(RegisterMyHandlers),
)

func NewMyRepository(db *gorm.DB) myrepo.Repository {
    return myrepo.New(db)
}

func NewMyUseCase(repo myrepo.Repository) myusecase.UseCase {
    return myusecase.New(repo)
}

func RegisterMyHandlers(params RouterParams, uc myusecase.UseCase) {
    handler.RegisterEndpoint(params.APIRouter, handler.NewHandler(uc))
}
```

### Register in main.go

```go
app := fx.New(
    bootstrap.Logger(),
    // Infrastructure
    bootstrap.Config,
    bootstrap.Database,
    // ...
    // Business Modules
    modules.Health,
    modules.MyModule,  // <- Add your new module
    // ...
)
```

## Benefits

- ✅ **No boilerplate** - Framework handles dependency wiring
- ✅ **Type-safe** - Compile-time validation of dependencies
- ✅ **Graceful shutdown** - All resources cleaned up properly
- ✅ **Testable** - Easy to inject mocks for testing
- ✅ **Modular** - Each domain is independent and self-contained
- ✅ **Clear separation** - Infrastructure vs Business Logic

## References

- [Uber FX Documentation](https://uber-go.github.io/fx/)

