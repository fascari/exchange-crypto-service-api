# Rate Limiter Implementation

This document describes the rate limiter implementation for the cryptocurrency exchange API.

## Features

### 1. IP-based Rate Limiting
- Limits requests per IP address
- Flexible configuration through `env.yaml`
- Automatic cleanup of inactive limiters
- Thread-safe for concurrent usage

## Configuration

### Configuration via `env.yaml`
Rate limiter configurations are loaded from the `env.yaml` file:

```yaml
# Rate Limiter Configuration
RATE_LIMITER_REQUESTS_PER_SECOND: 10      # Requests per second per IP
RATE_LIMITER_BURST_SIZE: 20               # Allowed burst size
RATE_LIMITER_CLEANUP_INTERVAL_MINUTES: 5  # Cleanup interval in minutes
```

### Current Configuration
The application is configured with:
- **10 requests per second** per IP (configurable)
- **Burst of 20 requests** (configurable)
- **Automatic cleanup** every 5 minutes (configurable)

### Default Values
If the configuration file is not found or is invalid, the default values are:
- `RequestsPerSecond`: 10
- `BurstSize`: 20
- `CleanupInterval`: 5 minutes

## How It Works

The rate limiter uses a token bucket algorithm to manage request rates. Each IP address has its own token bucket, which refills at a rate defined by `RequestsPerSecond` and allows bursts up to `BurstSize`.

### IP Headers
The rate limiter extracts the real client IP by checking:
1. `X-Forwarded-For` (proxies/load balancers)
2. `X-Real-IP` (alternative proxies)  
3. `RemoteAddr` (direct connection)

### Rate Limit Response
When the limit is exceeded:
```http
HTTP/1.1 429 Too Many Requests
Content-Type: application/json

{
  "error": "rate limit exceeded",
  "message": "too many requests"
}
```

## Customization

### Changing Configuration
To modify the limits, edit the `env.yaml`:

```yaml
# More restrictive configuration
RATE_LIMITER_REQUESTS_PER_SECOND: 5
RATE_LIMITER_BURST_SIZE: 10
RATE_LIMITER_CLEANUP_INTERVAL_MINUTES: 2

# More permissive configuration
RATE_LIMITER_REQUESTS_PER_SECOND: 20
RATE_LIMITER_BURST_SIZE: 50
RATE_LIMITER_CLEANUP_INTERVAL_MINUTES: 10
```