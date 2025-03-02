package health

import (
	"encoding/json"
	"net/http"
	"sync"
)

// Status represents the health status of the application
type Status string

const (
	// Up indicates the application is healthy
	Up Status = "UP"
	
	// Down indicates the application is unhealthy
	Down Status = "DOWN"
)

// Response represents the health check response
type Response struct {
	Status Status `json:"status"`
	Reason string `json:"reason,omitempty"`
}

var (
	// Default health status is UP
	currentStatus = Up
	
	// Reason for unhealthy status
	currentReason string
	
	// Mutex to protect concurrent access to health status
	mu sync.RWMutex
)

// SetHealthy sets the health status to UP and clears any reason
func SetHealthy() {
	mu.Lock()
	defer mu.Unlock()
	
	currentStatus = Up
	currentReason = ""
}


// SetUnhealthy sets the health status to DOWN with the given reason
func SetUnhealthy(reason string) {
	mu.Lock()
	defer mu.Unlock()
	
	currentStatus = Down
	currentReason = reason
}

// SetStatus sets the health status
func SetStatus(status Status) {
	mu.Lock()
	defer mu.Unlock()
	
	currentStatus = status
}

// SetReason sets the reason for the current health status
func SetReason(reason string) {
	mu.Lock()
	defer mu.Unlock()
	
	currentReason = reason
}

// GetStatus returns the current health status
func GetStatus() Status {
	mu.RLock()
	defer mu.RUnlock()
	
	return currentStatus
}

// GetReason returns the current reason for the health status
func GetReason() string {
	mu.RLock()
	defer mu.RUnlock()
	
	return currentReason
}

// HealthHandler represents a health check HTTP handler
type HealthHandler struct {
	useJSON bool
}

// NewHandler creates a new health check handler
func NewHandler() *HealthHandler {
	return &HealthHandler{
		useJSON: false,
	}
}

// WithJSON configures the handler to return JSON responses
func (h *HealthHandler) WithJSON(useJSON bool) *HealthHandler {
	h.useJSON = useJSON
	return h
}

func (h *HealthHandler)  getResponseStatusCodeAndBody() (statusCode int, body []byte) {
	mu.RLock()
	defer mu.RUnlock()
	
	h.useJSON = useJSON
	status := currentStatus
	reason := currentReason
}

// ServeHTTP implements the http.Handler interface
func (h *HealthHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	mu.RLock()
	status := currentStatus
	reason := currentReason
	mu.RUnlock()
	
	// Set appropriate status code
	if status == Down {
		w.WriteHeader(http.StatusServiceUnavailable)
	} else {
		w.WriteHeader(http.StatusOK)
	}
	
	// Create response
	response := Response{
		Status: status,
	}
	
	if status == Down && reason != "" {
		response.Reason = reason
	}
	
	// Return response in requested format
	if h.useJSON {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	} else {
		w.Header().Set("Content-Type", "text/plain")
		if status == Down && reason != "" {
			w.Write([]byte(string(status) + ": " + reason))
		} else {
			w.Write([]byte(string(status)))
		}
	}
}

// Handler returns a new health check HTTP handler
func Handler() *HealthHandler {
	return NewHandler()
}

