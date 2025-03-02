package main

import (
	"time"

	"github.com/andres-vara/health"
)

func main() {
    // Start with healthy status
    health.SetHealthy()

    // Start health check monitors
    go monitorServiceHealth()

    // ... rest of your service code ...
}

func monitorServiceHealth() {
    ticker := time.NewTicker(30 * time.Second)
    defer ticker.Stop()

    for range ticker.C {
        // Check critical dependencies
        if err := checkDatabaseConnection(); err != nil {
            health.SetUnhealthy("Database connection failed: " + err.Error())
            continue
        }

        // Check API dependencies
        if err := checkExternalAPIs(); err != nil {
            health.SetUnhealthy("External API check failed: " + err.Error())
            continue
        }

        // Check system resources
        if err := checkSystemResources(); err != nil {
            health.SetUnhealthy("System resource check failed: " + err.Error())
            continue
        }

        // If all checks pass, set as healthy
        health.SetHealthy()
    }
}

func checkDatabaseConnection() error {
    // Your database check logic
    return nil
}

func checkExternalAPIs() error {
    // Your API check logic
    return nil
}

func checkSystemResources() error {
    // Your resource check logic
    return nil
}