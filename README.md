# Health Package

This package provides a simple health check implementation for HTTP services. It supports:

- Basic health status reporting (UP/DOWN)
- Custom status messages/reasons
- Plain text and JSON response formats
- Integration with the `shttp` framework

## Features

- **Status Management**: Set application health status to UP or DOWN with optional reason messages
- **Response Format**: Choose between plain text and JSON responses 
- **Thread-safe**: All operations are protected by a mutex for concurrent access
- **Framework Integration**: Works with both standard `http.Handler` and the custom `shttp.Handler` pattern

## Integration with shttp

The health package now includes integration with the `shttp` framework via the `HealthHandler()` and `JSONHealthHandler()` functions that return a `health.Handler` type that matches the `shttp.Handler` interface:

```go
// Register a plain text health endpoint with an shttp router
server.GET("/health", health.HealthHandler())

// Register a JSON health endpoint with an shttp router
server.GET("/health/json", health.JSONHealthHandler())
```

The `health.Handler` type matches the `shttp.Handler` interface with these features:

- Context support
- Integration with request IDs for tracing
- Error handling

## Handler Types

The package defines these handler types:

```go
// Standard http.Handler interface implementation
func Handle() *healthHandler

// Handler compatible with shttp framework
func HealthHandler() Handler

// JSON-specific handler compatible with shttp framework
func JSONHealthHandler() Handler
```

## Usage Examples

### Standard HTTP Server

```go
import (
    "net/http"
    "github.com/yourusername/health"
)

func main() {
    // Configure the handler (optional, defaults to plain text)
    health.Handle().WithJSON(true)
    
    // Set initial status
    health.SetHealthy()
    
    // Register the endpoint
    http.Handle("/health", health.Handle())
    
    // Start the server
    http.ListenAndServe(":8080", nil)
}
```

### With shttp Framework

```go
import (
    "context"
    "github.com/yourusername/health"
    "github.com/yourusername/shttp"
    "github.com/yourusername/slogr"
)

func main() {
    // Create server
    ctx := context.Background()
    logger := slogr.New(os.Stdout, nil)
    config := shttp.DefaultConfig()
    config.Logger = logger
    
    server := shttp.New(ctx, config)
    
    // Add middleware
    server.Use(
        shttp.RequestIDMiddleware(),
        shttp.LoggingMiddleware(logger)
    )
    
    // Set initial health status
    health.SetHealthy()
    
    // Register standard health endpoint (plain text)
    server.GET("/health", health.HealthHandler())
    
    // Register JSON health endpoint (using the dedicated JSON handler)
    server.GET("/health/json", health.JSONHealthHandler())
    
    // Start the server
    server.Start()
}
```

## License

See the [LICENSE](LICENSE) file for details.
