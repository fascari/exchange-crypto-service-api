# Authentication (JWT)

The API uses JWT (JSON Web Tokens) for stateless authentication.

## How it works
1. **Token Generation**: Send user ID via path parameter to `/api/v1/tokens/generate/{user_id}` to get a JWT token
2. **Token Usage**: Include the token in the `Authorization: Bearer <token>` header for protected endpoints
3. **Token Validation**: The API automatically validates tokens and extracts user information

## JWT Token Structure
The JWT token contains the following claims:
- **user_id**: Unique user identifier from the database
- **username**: User identifier  
- **document_number**: User document number from the database
- **exp**: Token expiration timestamp
- **iat**: Token issued at timestamp
- **Signature**: HMAC-SHA256 signed for security

## Token Generation Request
**Path Parameter**: `user_id` (required) - included in the URL path

## Token Generation Response
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "expires_at": "2025-10-02T15:04:05Z"
}
```

## Example Usage
```bash
# Generate token (user must exist in database)
curl -X POST "http://localhost:8080/api/v1/tokens/generate/1"

# Use token in requests  
curl -X GET http://localhost:8080/api/v1/protected-endpoint \
  -H "Authorization: Bearer <your-jwt-token>"
```

## Requirements
- User must be registered via `/api/v1/users` endpoint before generating tokens
- Only existing users (found by user_id) can obtain JWT tokens
- All user data (username, document_number) is retrieved from the database
- user_id path parameter must be a valid positive integer

## Configuration
JWT settings can be configured in `env.yaml`:
- `JWT_SECRET`: Secret key for token signing
- `JWT_EXPIRATION_HOURS`: Token validity duration
