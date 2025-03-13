package health

// Example_shttp demonstrates how to use the health package with shttp
func Example_shttp() {
	/*
	This example shows how to use the health package with shttp.
	
	Note: This example won't compile directly as it references packages 
	that need to be imported in your actual project:
	- github.com/yourusername/shttp
	- github.com/yourusername/slogr
	
	The health.Handler type is compatible with shttp.Handler, making integration easy.

	Example usage:

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
		server := shttp.New(ctx, &shttp.Config{Logger: logger})
		
		// Register plain text health endpoint
		// health.Handler is compatible with shttp.Handler
		server.GET("/health", health.HealthHandler())
		
		// Register JSON health endpoint
		server.GET("/health/json", health.JSONHealthHandler())
		
		// Start server
		server.Start()
	}
	```
	*/
	
	// This is simplified example code
	
	// 1. Create context, logger, and server config
	// ctx := context.Background()
	// logger := slogr.New(os.Stdout, nil)
	// config := shttp.DefaultConfig()
	// config.Logger = logger
	
	// 2. Create server
	// server := shttp.New(ctx, config)
	
	// 3. Add middleware
	// server.Use(
	//    shttp.RequestIDMiddleware(),
	//    shttp.LoggingMiddleware(logger)
	// )
	
	// 4. Set initial health status
	// health.SetHealthy()
	
	// 5. Register health endpoints with the new Handler type
	// server.GET("/health", health.HealthHandler())
	
	// 6. Register JSON health endpoint 
	// server.GET("/health/json", health.JSONHealthHandler())
	
	// 7. Start the server
	// server.Start()
} 