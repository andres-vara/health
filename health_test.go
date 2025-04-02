package health

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHealthHandler(t *testing.T) {
	// Reset health status before each test
	SetHealthy()

	tests := []struct {
		name           string
		useJSON        bool
		setStatus      func()
		expectedStatus int
		expectedBody   string
		checkJSON      bool
	}{
		{
			name:           "Default UP status with plain text",
			useJSON:        false,
			setStatus:      func() { SetHealthy() },
			expectedStatus: http.StatusOK,
			expectedBody:   "UP: ",
		},
		{
			name:           "DOWN status with reason in plain text",
			useJSON:        false,
			setStatus:      func() { SetUnhealthy("Test reason") },
			expectedStatus: http.StatusServiceUnavailable,
			expectedBody:   "DOWN: Test reason",
		},
		{
			name:           "UP status with JSON",
			useJSON:        true,
			setStatus:      func() { SetHealthy() },
			expectedStatus: http.StatusOK,
			checkJSON:      true,
		},
		{
			name:           "DOWN status with reason in JSON",
			useJSON:        true,
			setStatus:      func() { SetUnhealthy("Test reason") },
			expectedStatus: http.StatusServiceUnavailable,
			checkJSON:      true,
		},
		{
			name:           "Set status and reason separately",
			useJSON:        true,
			setStatus:      func() { SetStatus(Down); SetReason("Manual reason") },
			expectedStatus: http.StatusServiceUnavailable,
			checkJSON:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set up the health status for this test
			tt.setStatus()

			// Create a request to pass to our handler
			req, err := http.NewRequest("GET", "/health", nil)
			if err != nil {
				t.Fatal(err)
			}

			// Create a ResponseRecorder to record the response
			rr := httptest.NewRecorder()

			// Create the handler with the appropriate JSON setting
			handler := Handle().WithJSON(tt.useJSON)

			// Serve the request
			handler.ServeHTTP(rr, req)

			// Check the status code
			if status := rr.Code; status != tt.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, tt.expectedStatus)
			}

			// Check the response body
			if tt.checkJSON {
				// For JSON responses, parse and check the structure
				var response responseBody
				body, _ := io.ReadAll(rr.Body)
				if err := json.Unmarshal(body, &response); err != nil {
					t.Errorf("Failed to parse JSON response: %v", err)
				}

				// Check status
				expectedStatus := "UP"
				if tt.expectedStatus != http.StatusOK {
					expectedStatus = "DOWN"
				}
				if response.Status != expectedStatus {
					t.Errorf("handler returned wrong status: got %v want %v",
						response.Status, expectedStatus)
				}

				// Check reason for DOWN status
				if expectedStatus == "DOWN" {
					if response.Reason == "" {
						t.Error("Expected reason to be set for DOWN status")
					}
				} else {
					if response.Reason != "" {
						t.Errorf("Expected no reason for UP status, got: %s", response.Reason)
					}
				}
			} else {
				// For plain text responses, check the exact body
				if body := rr.Body.String(); body != tt.expectedBody {
					t.Errorf("handler returned unexpected body: got %v want %v",
						body, tt.expectedBody)
				}
			}
		})
	}
}

func TestSHTTPHealthHandler(t *testing.T) {
	// Reset health status before each test
	SetHealthy()

	tests := []struct {
		name           string
		setStatus      func()
		expectedStatus int
		requestID      string
	}{
		{
			name:           "UP status with request ID",
			setStatus:      func() { SetHealthy() },
			expectedStatus: http.StatusOK,
			requestID:      "test-request-id",
		},
		{
			name:           "DOWN status with reason and request ID",
			setStatus:      func() { SetUnhealthy("Test reason") },
			expectedStatus: http.StatusServiceUnavailable,
			requestID:      "test-request-id-2",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set up the health status for this test
			tt.setStatus()

			// Create a request to pass to our handler
			req, err := http.NewRequest("GET", "/health", nil)
			if err != nil {
				t.Fatal(err)
			}

			// Create a ResponseRecorder to record the response
			rr := httptest.NewRecorder()

			// Create context with request ID
			ctx := context.WithValue(req.Context(), "request_id", tt.requestID)

			// Get the shttp handler
			handler := HealthHandler()

			// Call the handler with context
			err = handler(ctx, rr, req)
			if err != nil {
				t.Errorf("handler returned error: %v", err)
			}

			// Check the status code
			if status := rr.Code; status != tt.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, tt.expectedStatus)
			}

			// Check content type is JSON
			if contentType := rr.Header().Get("Content-Type"); contentType != "application/json" {
				t.Errorf("handler returned wrong content type: got %v want %v",
					contentType, "application/json")
			}

			// Check that request ID was forwarded to response header
			if requestID := rr.Header().Get("X-Request-ID"); requestID != tt.requestID {
				t.Errorf("handler did not forward request ID: got %v want %v",
					requestID, tt.requestID)
			}

			// Check the response body is valid JSON
			var response responseBody
			body, _ := io.ReadAll(rr.Body)
			if err := json.Unmarshal(body, &response); err != nil {
				t.Errorf("Failed to parse JSON response: %v", err)
			}

			// Check status
			expectedStatus := "UP"
			if tt.expectedStatus != http.StatusOK {
				expectedStatus = "DOWN"
			}
			if response.Status != expectedStatus {
				t.Errorf("handler returned wrong status: got %v want %v",
					response.Status, expectedStatus)
			}

			// Check reason for DOWN status
			if expectedStatus == "DOWN" {
				if response.Reason == "" {
					t.Error("Expected reason to be set for DOWN status")
				}
			} else {
				if response.Reason != "" {
					t.Errorf("Expected no reason for UP status, got: %s", response.Reason)
				}
			}
		})
	}
}

func TestSHTTPJSONHealthHandler(t *testing.T) {
	// Reset health status before each test
	SetHealthy()

	tests := []struct {
		name           string
		setStatus      func()
		expectedStatus int
		requestID      string
	}{
		{
			name:           "UP status with request ID",
			setStatus:      func() { SetHealthy() },
			expectedStatus: http.StatusOK,
			requestID:      "test-request-id",
		},
		{
			name:           "DOWN status with reason and request ID",
			setStatus:      func() { SetUnhealthy("Test reason") },
			expectedStatus: http.StatusServiceUnavailable,
			requestID:      "test-request-id-2",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set up the health status for this test
			tt.setStatus()

			// Create a request to pass to our handler
			req, err := http.NewRequest("GET", "/health", nil)
			if err != nil {
				t.Fatal(err)
			}

			// Create a ResponseRecorder to record the response
			rr := httptest.NewRecorder()

			// Create context with request ID
			ctx := context.WithValue(req.Context(), "request_id", tt.requestID)

			// Get the JSON shttp handler
			handler := JSONHealthHandler()

			// Call the handler with context
			err = handler(ctx, rr, req)
			if err != nil {
				t.Errorf("handler returned error: %v", err)
			}

			// Check the status code
			if status := rr.Code; status != tt.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, tt.expectedStatus)
			}

			// Check content type is JSON
			if contentType := rr.Header().Get("Content-Type"); contentType != "application/json" {
				t.Errorf("handler returned wrong content type: got %v want %v",
					contentType, "application/json")
			}

			// Check that request ID was forwarded to response header
			if requestID := rr.Header().Get("X-Request-ID"); requestID != tt.requestID {
				t.Errorf("handler did not forward request ID: got %v want %v",
					requestID, tt.requestID)
			}

			// Check the response body is valid JSON
			var response responseBody
			body, _ := io.ReadAll(rr.Body)
			if err := json.Unmarshal(body, &response); err != nil {
				t.Errorf("Failed to parse JSON response: %v", err)
			}

			// Check status
			expectedStatus := "UP"
			if tt.expectedStatus != http.StatusOK {
				expectedStatus = "DOWN"
			}
			if response.Status != expectedStatus {
				t.Errorf("handler returned wrong status: got %v want %v",
					response.Status, expectedStatus)
			}

			// Check reason for DOWN status
			if expectedStatus == "DOWN" {
				if response.Reason == "" {
					t.Error("Expected reason to be set for DOWN status")
				}
			} else {
				if response.Reason != "" {
					t.Errorf("Expected no reason for UP status, got: %s", response.Reason)
				}
			}
		})
	}
}

func TestStatusManagement(t *testing.T) {
	// Test SetHealthy
	SetHealthy()
	if status := GetStatus(); status != Up {
		t.Errorf("SetHealthy failed: got %v want %v", status, Up)
	}
	if reason := GetReason(); reason != "" {
		t.Errorf("SetHealthy should clear reason: got %v want %v", reason, "")
	}

	// Test SetUnhealthy
	SetUnhealthy("Test reason")
	if status := GetStatus(); status != Down {
		t.Errorf("SetUnhealthy failed: got %v want %v", status, Down)
	}
	if reason := GetReason(); reason != "Test reason" {
		t.Errorf("SetUnhealthy reason mismatch: got %v want %v", reason, "Test reason")
	}

	// Test SetStatus
	SetStatus(Up)
	if status := GetStatus(); status != Up {
		t.Errorf("SetStatus failed: got %v want %v", status, Up)
	}

	// Test SetReason
	SetReason("New reason")
	if reason := GetReason(); reason != "New reason" {
		t.Errorf("SetReason failed: got %v want %v", reason, "New reason")
	}
}

func TestConcurrentAccess(t *testing.T) {
	// Reset health status
	SetHealthy()

	// Create a done channel to signal completion
	done := make(chan bool)

	// Start multiple goroutines to access and modify health status
	for i := 0; i < 10; i++ {
		go func(id int) {
			// Toggle health status
			if id%2 == 0 {
				SetHealthy()
			} else {
				SetUnhealthy("Reason from goroutine " + string(rune(id+'0')))
			}

			// Read health status
			_ = GetStatus()
			_ = GetReason()

			done <- true
		}(i)
	}

	// Wait for all goroutines to complete
	for i := 0; i < 10; i++ {
		<-done
	}

	// If we got here without deadlock or panic, the test passes
} 