package test

import (
	"net/http"
	"testing"
)

// TestCallEndpoint tests the CallEndpoint function
func TestCallEndpoint(t *testing.T) {
	// Test with a simple handler
	handler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Hello, World!"))
	}

	// Test with default options
	body, resp, err := CallEndpoint("GET", handler, NewRequestOptions{})
	if err != nil {
		t.Fatalf("CallEndpoint failed: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, resp.StatusCode)
	}
	if body != "Hello, World!" {
		t.Errorf("Expected body %q, got %q", "Hello, World!", body)
	}

	// Test with custom body
	body, resp, err = CallEndpoint("POST", handler, NewRequestOptions{
		Body: "Custom body",
	})
	if err != nil {
		t.Fatalf("CallEndpoint failed: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, resp.StatusCode)
	}

	// Test with custom headers
	body, resp, err = CallEndpoint("GET", handler, NewRequestOptions{
		Headers: map[string]string{
			"X-Test-Header": "test-value",
		},
	})
	if err != nil {
		t.Fatalf("CallEndpoint failed: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, resp.StatusCode)
	}
}

// TestCallStringEndpoint tests the CallStringEndpoint function
func TestCallStringEndpoint(t *testing.T) {
	// Test with a simple handler
	handler := func(w http.ResponseWriter, r *http.Request) string {
		return "Hello, String World!"
	}

	// Test with default options
	body, resp, err := CallStringEndpoint("GET", handler, NewRequestOptions{})
	if err != nil {
		t.Fatalf("CallStringEndpoint failed: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, resp.StatusCode)
	}
	if body != "Hello, String World!" {
		t.Errorf("Expected body %q, got %q", "Hello, String World!", body)
	}

	// Test with custom body
	body, resp, err = CallStringEndpoint("POST", handler, NewRequestOptions{
		Body: "Custom body",
	})
	if err != nil {
		t.Fatalf("CallStringEndpoint failed: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, resp.StatusCode)
	}

	// Test with custom headers
	body, resp, err = CallStringEndpoint("GET", handler, NewRequestOptions{
		Headers: map[string]string{
			"X-Test-Header": "test-value",
		},
	})
	if err != nil {
		t.Fatalf("CallStringEndpoint failed: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, resp.StatusCode)
	}
}

// TestStringHandler tests the StringHandler type
func TestStringHandler(t *testing.T) {
	// Create a StringHandler
	handler := StringHandler(func(w http.ResponseWriter, r *http.Request) string {
		return "Hello from StringHandler!"
	})

	// Create a test request
	req := NewTestHTTPRequest("GET", "/")
	recorder := req.Execute(handler)

	// Check the response
	if recorder.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, recorder.Code)
	}
	if recorder.Body.String() != "Hello from StringHandler!" {
		t.Errorf("Expected body %q, got %q", "Hello from StringHandler!", recorder.Body.String())
	}
}
