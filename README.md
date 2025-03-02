# Health

A simple Go library for implementing health check endpoints in web applications.

## Installation

```bash
go get github.com/andres-vara/health
```

## Usage

### Basic Usage

```go
package main

import (
    "net/http"
    
    "github.com/andres-vara/health"
)

func main() {
    // Create a new router
    router := http.NewServeMux()
    
    // Add health endpoint
    router.Handle("/health", health.Handler())
    
    // Start server
    server := &http.Server{
        Addr:    ":8080",
        Handler: router,
    }
    
    server.ListenAndServe()
}
```

### JSON Response

By default, the health endpoint returns plain text responses. To enable JSON responses:

```go
router.Handle("/health", health.Handler().WithJSON(true))
```

### Setting Health Status

The health status is "UP" by default. To mark the service as unhealthy:

```go
// Set status to DOWN with a reason
health.SetUnhealthy("Database connection failed")

// Set status back to UP and clear the reason
health.SetHealthy()
```

You can also set the status and reason separately:

```go
// Set status to DOWN
health.SetStatus(health.Down)

// Set a reason
health.SetReason("Redis cache unavailable")
```

### Response Format

#### Plain Text (default)
- When healthy: `UP`
- When unhealthy: `DOWN: reason`

#### JSON (when enabled)
- When healthy: `{"status": "UP"}`
- When unhealthy: `{"status": "DOWN", "reason": "reason"}`

### HTTP Status Codes

- `200 OK` when the status is UP
- `503 Service Unavailable` when the status is DOWN

## License

See the [LICENSE](LICENSE) file for details.
