package health

import (
	"encoding/json"
	"net/http"
	"sync"
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

func (h *healthHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	statusCode, body, useJSON := h.getStatus()

	if useJSON {
		w.Header().Set("Content-Type", "application/json")
	}

	w.WriteHeader(statusCode)

	_, _ = w.Write(body)
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
