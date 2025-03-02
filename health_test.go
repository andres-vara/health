package health

import (
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
			expectedBody:   "UP",
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
			handler := NewHandler()
			if tt.useJSON {
				handler = handler.WithJSON(true)
			}

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
				var response Response
				body, _ := io.ReadAll(rr.Body)
				if err := json.Unmarshal(body, &response); err != nil {
					t.Errorf("Failed to parse JSON response: %v", err)
				}

				// Check status
				expectedStatus := Up
				if tt.expectedStatus != http.StatusOK {
					expectedStatus = Down
				}
				if response.Status != expectedStatus {
					t.Errorf("handler returned wrong status: got %v want %v",
						response.Status, expectedStatus)
				}

				// Check reason for DOWN status
				if expectedStatus == Down {
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