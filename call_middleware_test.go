package test_test

import (
	"net/http"
	"testing"

	"github.com/dracory/test"
)

// TestCallMiddleware tests the CallMiddleware function
func TestCallMiddleware(t *testing.T) {
	// Create a middleware
	middleware := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("X-Test-Middleware", "true")
			next.ServeHTTP(w, r)
		})
	}

	// Create a handler
	handler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Middleware test"))
	}

	// Test with default options
	body, resp, err := test.CallMiddleware("GET", middleware, handler, test.NewRequestOptions{})
	if err != nil {
		t.Fatalf("CallMiddleware failed: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, resp.StatusCode)
	}
	if body != "Middleware test" {
		t.Errorf("Expected body %q, got %q", "Middleware test", body)
	}
	if resp.Header.Get("X-Test-Middleware") != "true" {
		t.Errorf("Expected X-Test-Middleware header to be set")
	}

	// Test with custom options
	body, resp, err = test.CallMiddleware("POST", middleware, handler, test.NewRequestOptions{
		Body: "Custom body",
		Headers: map[string]string{
			"X-Test-Header": "test-value",
		},
	})
	if err != nil {
		t.Fatalf("CallMiddleware failed: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, resp.StatusCode)
	}
	if resp.Header.Get("X-Test-Middleware") != "true" {
		t.Errorf("Expected X-Test-Middleware header to be set")
	}
	if body != "Middleware test" {
		t.Errorf("Expected body %q, got %q", "Middleware test", body)
	}
}
