package health

import (
	"context"
	"encoding/json"
	"net/http"
	"sync"

	"github.com/andres-vara/shttp"
)

// We'll define our own type for context keys to avoid dependency on shttp
type ContextKey string

const (
	// These match the keys in shttp package
	LoggerKey   ContextKey = "logger"
	RequestIDKey ContextKey = "request_id"
)

type Status string

var (
	Up Status = "UP"
	Down Status = "DOWN"
	handler  = &healthHandler{
		status: Up,
		useJSON: false,
	}
)

type responseBody struct {
	Status string `json:"status"`
	Reason string `json:"reason,omitempty"`
}

type healthHandler struct {
	status Status
	reason string

	useJSON bool
	mutex sync.RWMutex
}

// ServeHTTP implements the http.Handler interface for standard HTTP servers
func (h *healthHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	statusCode, body, useJSON := h.getStatus()

	if useJSON {
		w.Header().Set("Content-Type", "application/json")
	}

	w.WriteHeader(statusCode)

	_, _ = w.Write(body)
}

// HealthHandler returns a handler compatible with shttp.Handler interface
// for use with the shttp package. This uses the default format (plain text or JSON)
// based on the current settings of the health handler.
func HealthHandler() shttp.Handler {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
		// Get status information
		statusCode, body, useJSON := handler.getStatus()

		// Set appropriate content type
		if useJSON {
			w.Header().Set("Content-Type", "application/json")
		}

		// Forward any request ID from context to response headers for traceability
		if requestID, ok := ctx.Value("request_id").(string); ok && requestID != "" {
			w.Header().Set("X-Request-ID", requestID)
		}

		// Set status code and write response
		w.WriteHeader(statusCode)
		_, _ = w.Write(body)
		
		return nil
	}
}

// JSONHealthHandler returns a handler that always returns JSON responses,
// regardless of the current handler configuration.
func JSONHealthHandler() shttp.Handler {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
		// Get the current status but force JSON format
		handler.mutex.RLock()
		status := handler.status
		reason := handler.reason
		handler.mutex.RUnlock()
		
		// Create JSON response
		body, _ := json.Marshal(responseBody{
			Status: string(status),
			Reason: reason,
		})
		
		// Set appropriate headers
		w.Header().Set("Content-Type", "application/json")
		
		// Forward any request ID from context
		if requestID, ok := ctx.Value("request_id").(string); ok && requestID != "" {
			w.Header().Set("X-Request-ID", requestID)
		}
		
		// Set status code
		statusCode := http.StatusOK
		if status == Down {
			statusCode = http.StatusServiceUnavailable
		}
		
		w.WriteHeader(statusCode)
		_, _ = w.Write(body)
		
		return nil
	}
}

func (h *healthHandler) GetResponseStatusCodeAndBody() (int, []byte) {
	statusCode, body, _ := h.getStatus()
	return statusCode, body
}

func (h *healthHandler) getStatus() (int, []byte, bool) {
	var status Status
	var reason string
	var body []byte
	var useJSON bool
	var statusCode int

	h.mutex.RLock()
	status = h.status
	reason = h.reason
	useJSON = h.useJSON
	h.mutex.RUnlock()

	if useJSON {
		body, _ = json.Marshal(responseBody{
			Status: string(status),
			Reason: reason,
		})
	} else {
		body = []byte(string(status) + ": " + reason)
	}

	if status == Up {
		statusCode = http.StatusOK
	} else {
		statusCode = http.StatusServiceUnavailable
	}

	return statusCode, body, useJSON
}

func Handle() *healthHandler {
	return handler
}

func GetStatus() Status {
	handler.mutex.RLock()
	defer handler.mutex.RUnlock()

	return handler.status
}

func SetStatus(status Status) {
	handler.mutex.Lock()
	defer handler.mutex.Unlock()

	handler.status = status
}

func SetReason(reason string) {
	handler.mutex.Lock()
	defer handler.mutex.Unlock()

	handler.reason = reason
}

func GetReason() string {
	handler.mutex.RLock()
	defer handler.mutex.RUnlock()

	return handler.reason
}

func SetHealthy() {
	SetStatus(Up)
	SetReason("")
}

func SetUnhealthy(reason string) {
	SetStatus(Down)
	SetReason(reason)
}

func (h *healthHandler) WithJSON(v bool) *healthHandler {
	h.useJSON = v
	return h
}
