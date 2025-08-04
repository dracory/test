package test

import (
	"net/http"
	"testing"
)

// TestStringResponse tests the StringResponse function
func TestStringResponse(t *testing.T) {
	// Test with no content type set
	recorder := NewTestHTTPRequest("GET", "/").Execute(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		StringResponse(w, r, "Test response")
	}))

	if recorder.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, recorder.Code)
	}
	if recorder.Body.String() != "Test response" {
		t.Errorf("Expected body %q, got %q", "Test response", recorder.Body.String())
	}
	contentType := recorder.Header().Get("Content-Type")
	if contentType != "text/html; charset=utf-8" {
		t.Errorf("Expected Content-Type %q, got %q", "text/html; charset=utf-8", contentType)
	}

	// Test with content type already set
	recorder = NewTestHTTPRequest("GET", "/").Execute(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		StringResponse(w, r, "Test response")
	}))

	contentType = recorder.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("Expected Content-Type %q, got %q", "application/json", contentType)
	}
}
