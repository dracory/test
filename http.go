package test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
)

// TestHTTPRequest represents a test HTTP request
type TestHTTPRequest struct {
	Method  string
	Path    string
	Body    io.Reader
	Headers map[string]string
}

// NewTestHTTPRequest creates a new test HTTP request
func NewTestHTTPRequest(method, path string) *TestHTTPRequest {
	return &TestHTTPRequest{
		Method:  method,
		Path:    path,
		Headers: make(map[string]string),
	}
}

// WithBody sets the request body
func (r *TestHTTPRequest) WithBody(body string) *TestHTTPRequest {
	r.Body = strings.NewReader(body)
	return r
}

// WithHeader adds a header to the request
func (r *TestHTTPRequest) WithHeader(key, value string) *TestHTTPRequest {
	r.Headers[key] = value
	return r
}

// WithJSONBody sets the request body as JSON and adds the appropriate content type header
func (r *TestHTTPRequest) WithJSONBody(jsonBody string) *TestHTTPRequest {
	r.Body = strings.NewReader(jsonBody)
	r.Headers["Content-Type"] = "application/json"
	return r
}

// WithFormBody sets the request body as form data and adds the appropriate content type header
func (r *TestHTTPRequest) WithFormBody(formBody string) *TestHTTPRequest {
	r.Body = strings.NewReader(formBody)
	r.Headers["Content-Type"] = "application/x-www-form-urlencoded"
	return r
}

// Execute executes the request against the provided handler
func (r *TestHTTPRequest) Execute(handler http.Handler) *httptest.ResponseRecorder {
	req := httptest.NewRequest(r.Method, r.Path, r.Body)

	// Set headers
	for key, value := range r.Headers {
		req.Header.Set(key, value)
	}

	// Create a response recorder
	recorder := httptest.NewRecorder()

	// Serve the request
	handler.ServeHTTP(recorder, req)

	return recorder
}

// TestHTTPServer is a wrapper around httptest.Server for testing HTTP servers
type TestHTTPServer struct {
	Server *httptest.Server
}

// NewTestHTTPServer creates a new test HTTP server with the provided handler
func NewTestHTTPServer(handler http.Handler) *TestHTTPServer {
	server := httptest.NewServer(handler)
	return &TestHTTPServer{
		Server: server,
	}
}

// URL returns the base URL of the test server
func (s *TestHTTPServer) URL() string {
	return s.Server.URL
}

// Close closes the test server
func (s *TestHTTPServer) Close() {
	s.Server.Close()
}

// Get performs a GET request to the test server
func (s *TestHTTPServer) Get(path string) (*http.Response, error) {
	return http.Get(s.URL() + path)
}

// Post performs a POST request to the test server
func (s *TestHTTPServer) Post(path string, contentType string, body io.Reader) (*http.Response, error) {
	return http.Post(s.URL()+path, contentType, body)
}

// Do performs a custom request to the test server
func (s *TestHTTPServer) Do(req *http.Request) (*http.Response, error) {
	// Update the URL to point to the test server
	req.URL.Scheme = "http"
	req.URL.Host = strings.TrimPrefix(s.URL(), "http://")

	// Perform the request
	client := &http.Client{}
	return client.Do(req)
}
