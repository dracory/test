package test_test

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/dracory/test"
)

// TestNewRequest tests the NewRequest function
func TestNewRequest(t *testing.T) {
	// Test with minimal options
	req, err := test.NewRequest("GET", "/test", test.NewRequestOptions{})
	if err != nil {
		t.Fatalf("NewRequest failed: %v", err)
	}
	if req.Method != "GET" {
		t.Errorf("Expected method %q, got %q", "GET", req.Method)
	}
	if req.URL.Path != "/test" {
		t.Errorf("Expected path %q, got %q", "/test", req.URL.Path)
	}
	if req.RequestURI != "/test" {
		t.Errorf("Expected RequestURI %q, got %q", "/test", req.RequestURI)
	}

	// Test with empty URL (should default to "/")
	req, err = test.NewRequest("GET", "", test.NewRequestOptions{})
	if err != nil {
		t.Fatalf("NewRequest failed: %v", err)
	}
	if req.URL.Path != "/" {
		t.Errorf("Expected path %q, got %q", "/", req.URL.Path)
	}

	// Test with body
	req, err = test.NewRequest("POST", "/test", test.NewRequestOptions{
		Body: "Test body",
	})
	if err != nil {
		t.Fatalf("NewRequest failed: %v", err)
	}
	body, _ := io.ReadAll(req.Body)
	if string(body) != "Test body" {
		t.Errorf("Expected body %q, got %q", "Test body", string(body))
	}

	// Test with JSON data
	type TestData struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}
	jsonData := TestData{Name: "John", Age: 30}
	req, err = test.NewRequest("POST", "/test", test.NewRequestOptions{
		JSONData: jsonData,
	})
	if err != nil {
		t.Fatalf("NewRequest failed: %v", err)
	}
	if req.Header.Get("Content-Type") != "application/json" {
		t.Errorf("Expected Content-Type %q, got %q", "application/json", req.Header.Get("Content-Type"))
	}

	// Test with form values
	formValues := url.Values{}
	formValues.Set("name", "John")
	formValues.Set("age", "30")
	req, err = test.NewRequest("POST", "/test", test.NewRequestOptions{
		FormValues: formValues,
	})
	if err != nil {
		t.Fatalf("NewRequest failed: %v", err)
	}
	if req.Header.Get("Content-Type") != "application/x-www-form-urlencoded" {
		t.Errorf("Expected Content-Type %q, got %q", "application/x-www-form-urlencoded", req.Header.Get("Content-Type"))
	}

	// Test with query parameters
	queryParams := url.Values{}
	queryParams.Set("q", "search")
	queryParams.Set("page", "1")
	req, err = test.NewRequest("GET", "/test", test.NewRequestOptions{
		QueryParams: queryParams,
	})
	if err != nil {
		t.Fatalf("NewRequest failed: %v", err)
	}
	expectedQuery := "page=1&q=search"
	if req.URL.RawQuery != expectedQuery {
		t.Errorf("Expected query %q, got %q", expectedQuery, req.URL.RawQuery)
	}

	// Test with headers
	headers := map[string]string{
		"X-Test-Header": "test",
		"Authorization": "Bearer token",
	}
	req, err = test.NewRequest("GET", "/test", test.NewRequestOptions{
		Headers: headers,
	})
	if err != nil {
		t.Fatalf("NewRequest failed: %v", err)
	}
	if req.Header.Get("X-Test-Header") != "test" {
		t.Errorf("Expected header X-Test-Header %q, got %q", "test", req.Header.Get("X-Test-Header"))
	}
	if req.Header.Get("Authorization") != "Bearer token" {
		t.Errorf("Expected header Authorization %q, got %q", "Bearer token", req.Header.Get("Authorization"))
	}

	// Test with context values
	type contextKey string
	const userKey contextKey = "user"
	contextValues := map[any]any{
		userKey: "testuser",
	}
	req, err = test.NewRequest("GET", "/test", test.NewRequestOptions{
		Context: contextValues,
	})
	if err != nil {
		t.Fatalf("NewRequest failed: %v", err)
	}
	if req.Context().Value(userKey) != "testuser" {
		t.Errorf("Expected context value %q, got %q", "testuser", req.Context().Value(userKey))
	}

	// Test with explicit content type
	req, err = test.NewRequest("POST", "/test", test.NewRequestOptions{
		Body:        "Test body",
		ContentType: "text/plain",
	})
	if err != nil {
		t.Fatalf("NewRequest failed: %v", err)
	}
	if req.Header.Get("Content-Type") != "text/plain" {
		t.Errorf("Expected Content-Type %q, got %q", "text/plain", req.Header.Get("Content-Type"))
	}

	// Test with deprecated GetValues (should set QueryParams)
	getValues := url.Values{}
	getValues.Set("q", "legacy")
	req, err = test.NewRequest("GET", "/test", test.NewRequestOptions{
		GetValues: getValues,
	})
	if err != nil {
		t.Fatalf("NewRequest failed: %v", err)
	}
	if req.URL.Query().Get("q") != "legacy" {
		t.Errorf("Expected query parameter q=%q, got %q", "legacy", req.URL.Query().Get("q"))
	}

	// Test with deprecated PostValues (should set FormValues)
	postValues := url.Values{}
	postValues.Set("name", "legacy")
	req, err = test.NewRequest("POST", "/test", test.NewRequestOptions{
		PostValues: postValues,
	})
	if err != nil {
		t.Fatalf("NewRequest failed: %v", err)
	}
	if req.Header.Get("Content-Type") != "application/x-www-form-urlencoded" {
		t.Errorf("Expected Content-Type %q, got %q", "application/x-www-form-urlencoded", req.Header.Get("Content-Type"))
	}
}

// TestTestHTTPRequest tests the TestHTTPRequest type and its methods
func TestTestHTTPRequest(t *testing.T) {
	// Test NewTestHTTPRequest
	req := test.NewTestHTTPRequest("GET", "/test")
	if req.Method != "GET" {
		t.Errorf("Expected method %q, got %q", "GET", req.Method)
	}
	if req.Path != "/test" {
		t.Errorf("Expected path %q, got %q", "/test", req.Path)
	}

	// Test WithBody
	req = req.WithBody("Test body")
	body, _ := io.ReadAll(req.Body)
	if string(body) != "Test body" {
		t.Errorf("Expected body %q, got %q", "Test body", string(body))
	}

	// Test WithHeader
	req = req.WithHeader("X-Test", "value")
	if req.Headers["X-Test"] != "value" {
		t.Errorf("Expected header X-Test %q, got %q", "value", req.Headers["X-Test"])
	}

	// Test WithJSONBody
	req = req.WithJSONBody(`{"name":"John"}`)
	body, _ = io.ReadAll(req.Body)
	if string(body) != `{"name":"John"}` {
		t.Errorf("Expected JSON body %q, got %q", `{"name":"John"}`, string(body))
	}
	if req.Headers["Content-Type"] != "application/json" {
		t.Errorf("Expected Content-Type %q, got %q", "application/json", req.Headers["Content-Type"])
	}

	// Test WithFormBody
	req = req.WithFormBody("name=John")
	body, _ = io.ReadAll(req.Body)
	if string(body) != "name=John" {
		t.Errorf("Expected form body %q, got %q", "name=John", string(body))
	}
	if req.Headers["Content-Type"] != "application/x-www-form-urlencoded" {
		t.Errorf("Expected Content-Type %q, got %q", "application/x-www-form-urlencoded", req.Headers["Content-Type"])
	}

	// Test Execute
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Test response"))
	})
	recorder := req.Execute(handler)
	if recorder.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, recorder.Code)
	}
	if recorder.Body.String() != "Test response" {
		t.Errorf("Expected body %q, got %q", "Test response", recorder.Body.String())
	}
}

// TestTestHTTPServer tests the TestHTTPServer type and its methods
func TestTestHTTPServer(t *testing.T) {
	// Create a test handler
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			body, _ := io.ReadAll(r.Body)
			w.Write([]byte(fmt.Sprintf("Received: %s", string(body))))
			return
		}
		w.Write([]byte("Hello from test server"))
	})

	// Create a test server
	server := test.NewTestHTTPServer(handler)
	defer server.Close()

	// Test URL
	url := server.URL()
	if !strings.HasPrefix(url, "http://") {
		t.Errorf("Expected URL to start with http://, got %q", url)
	}

	// Test Get
	resp, err := server.Get("/")
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	if string(body) != "Hello from test server" {
		t.Errorf("Expected body %q, got %q", "Hello from test server", string(body))
	}

	// Test Post
	resp, err = server.Post("/", "text/plain", strings.NewReader("Test post"))
	if err != nil {
		t.Fatalf("Post failed: %v", err)
	}
	defer resp.Body.Close()
	body, _ = io.ReadAll(resp.Body)
	if string(body) != "Received: Test post" {
		t.Errorf("Expected body %q, got %q", "Received: Test post", string(body))
	}

	// Test Do with custom request
	req, _ := http.NewRequest("GET", "/custom", nil)
	resp, err = server.Do(req)
	if err != nil {
		t.Fatalf("Do failed: %v", err)
	}
	defer resp.Body.Close()
	body, _ = io.ReadAll(resp.Body)
	if string(body) != "Hello from test server" {
		t.Errorf("Expected body %q, got %q", "Hello from test server", string(body))
	}
}
