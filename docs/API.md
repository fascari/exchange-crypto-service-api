# API Documentation

## Main Endpoints
- **POST /api/v1/tokens/generate**: generate JWT authentication token
- **POST /api/v1/users**: create user
- **GET /api/v1/users/{userId}/balance**: check balance
- **POST /api/v1/accounts**: create account
- **POST /api/v1/transactions**: create transaction (deposit/withdrawal)
- **GET /api/v1/transactions/daily?date=YYYY-MM-DD**: get daily transactions
- **GET /health**: health check

## Postman Collection

The project includes a comprehensive Postman collection with example requests for all API endpoints.
- Import `exchange-crypto.postman_collection.json` file from the project root

## Request/Response Examples

### Authentication
```bash
# Generate JWT token
curl -X POST http://localhost:8080/api/v1/tokens/generate \
  -H "Content-Type: application/json" \
  -d '{"username": "john", "email": "john@example.com", "password": "secret"}'
```

### User Operations
```bash
# Create user
curl -X POST http://localhost:8080/api/v1/users \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <token>" \
  -d '{"name": "John Doe", "email": "john@example.com"}'

# Get user balance
curl -X GET http://localhost:8080/api/v1/users/1/balance \
  -H "Authorization: Bearer <token>"
```

### Transaction Operations
```bash
# Create transaction
curl -X POST http://localhost:8080/api/v1/transactions \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <token>" \
  -d '{"type": "deposit", "amount": 100.50, "account_id": 1}'

# Get daily transactions
curl -X GET "http://localhost:8080/api/v1/transactions/daily?date=2024-01-15" \
  -H "Authorization: Bearer <token>"
```
