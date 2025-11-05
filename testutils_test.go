package test

import (
	"database/sql"
	"net/http"
	"os"
	"testing"

	_ "modernc.org/sqlite"
)

func TestTestKey(t *testing.T) {
	// Test that the TestKey function generates a consistent hash
	key1 := TestKey("sqlite", "", "", "test.db", "", "")
	key2 := TestKey("sqlite", "", "", "test.db", "", "")

	if key1 != key2 {
		t.Errorf("TestKey should generate consistent hashes for the same input")
	}

	// Test that different inputs generate different hashes
	key3 := TestKey("mysql", "localhost", "3306", "testdb", "root", "password")

	if key1 == key3 {
		t.Errorf("TestKey should generate different hashes for different inputs")
	}
}

func TestTestConfig(t *testing.T) {
	// Get the original value of an environment variable to restore later
	originalAppName := os.Getenv("APP_NAME")
	defer os.Setenv("APP_NAME", originalAppName)

	// Create a test configuration
	config := DefaultTestConfig()
	config.AppName = "Test App Name"

	// Set up the test environment
	SetupTestEnvironment(config)

	// Check that the environment variables were set
	if os.Getenv("APP_NAME") != "Test App Name" {
		t.Errorf("Expected APP_NAME to be 'Test App Name', got '%s'", os.Getenv("APP_NAME"))
	}

	if os.Getenv("APP_ENV") != EnvTesting {
		t.Errorf("Expected APP_ENV to be '%s', got '%s'", EnvTesting, os.Getenv("APP_ENV"))
	}

	// Clean up the test environment
	CleanupTestEnvironment(config)

	// Check that the environment variables were unset
	if os.Getenv("APP_NAME") != "" {
		t.Errorf("Expected APP_NAME to be unset, got '%s'", os.Getenv("APP_NAME"))
	}
}

func TestTestDB(t *testing.T) {
	// Create a test database
	db, err := NewTestDB(nil)
	if err != nil {
		t.Fatalf("Failed to create test database: %v", err)
	}
	defer CloseTestDB(db)

	// Create a test table
	err = CreateTestTable(db, "users", "id INTEGER PRIMARY KEY, name TEXT, email TEXT")
	if err != nil {
		t.Fatalf("Failed to create test table: %v", err)
	}

	// Insert test data
	err = ExecuteSQLWithArgs(db, "INSERT INTO users (name, email) VALUES (?, ?)", "Test User", "test@example.com")
	if err != nil {
		t.Fatalf("Failed to insert test data: %v", err)
	}

	// Query the test data
	var name string
	var email string
	err = db.QueryRow("SELECT name, email FROM users WHERE name = ?", "Test User").Scan(&name, &email)
	if err != nil {
		t.Fatalf("Failed to query test data: %v", err)
	}

	// Check the results
	if name != "Test User" {
		t.Errorf("Expected name to be 'Test User', got '%s'", name)
	}

	if email != "test@example.com" {
		t.Errorf("Expected email to be 'test@example.com', got '%s'", email)
	}

	// Drop the test table
	err = DropTestTable(db, "users")
	if err != nil {
		t.Fatalf("Failed to drop test table: %v", err)
	}

	// Verify the table was dropped
	_, err = db.Query("SELECT * FROM users")
	if err == nil {
		t.Errorf("Expected an error when querying a dropped table")
	}
}

func TestTestHTTP(t *testing.T) {
	// Create a test handler
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" && r.URL.Path == "/hello" {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("Hello, world!"))
			return
		}

		if r.Method == "POST" && r.URL.Path == "/echo" {
			// Read the request body
			buf := make([]byte, r.ContentLength)
			r.Body.Read(buf)

			// Echo the request body
			w.WriteHeader(http.StatusOK)
			w.Write(buf)
			return
		}

		w.WriteHeader(http.StatusNotFound)
	})

	// Test a GET request
	getReq := NewTestHTTPRequest("GET", "/hello")
	getResp := getReq.Execute(handler)

	if getResp.Code != http.StatusOK {
		t.Errorf("Expected status code %d for GET /hello, got %d", http.StatusOK, getResp.Code)
	}

	if getResp.Body.String() != "Hello, world!" {
		t.Errorf("Expected body %q for GET /hello, got %q", "Hello, world!", getResp.Body.String())
	}

	// Test a POST request with JSON body
	postReq := NewTestHTTPRequest("POST", "/echo").WithJSONBody(`{"message":"Hello"}`)
	postResp := postReq.Execute(handler)

	if postResp.Code != http.StatusOK {
		t.Errorf("Expected status code %d for POST /echo, got %d", http.StatusOK, postResp.Code)
	}

	if postResp.Body.String() != `{"message":"Hello"}` {
		t.Errorf("Expected body %q for POST /echo, got %q", `{"message":"Hello"}`, postResp.Body.String())
	}

	// Test a request to a non-existent endpoint
	notFoundReq := NewTestHTTPRequest("GET", "/not-found")
	notFoundResp := notFoundReq.Execute(handler)

	if notFoundResp.Code != http.StatusNotFound {
		t.Errorf("Expected status code %d for GET /not-found, got %d", http.StatusNotFound, notFoundResp.Code)
	}
}

// Example of a more complex test that combines multiple utilities
func TestIntegration(t *testing.T) {
	// Set up the test environment
	config := DefaultTestConfig()
	SetupTestEnvironment(config)
	defer CleanupTestEnvironment(config)

	// Create a test database
	db, err := NewTestDB(nil)
	if err != nil {
		t.Fatalf("Failed to create test database: %v", err)
	}
	defer CloseTestDB(db)

	// Create a test table
	err = CreateTestTable(db, "users", "id INTEGER PRIMARY KEY, name TEXT, email TEXT")
	if err != nil {
		t.Fatalf("Failed to create test table: %v", err)
	}

	// Create a handler that uses the database
	handler := createTestHandler(db)

	// Create a test server
	server := NewTestHTTPServer(handler)
	defer server.Close()

	// Test the handler
	req := NewTestHTTPRequest("POST", "/users").
		WithJSONBody(`{"name":"Test User","email":"test@example.com"}`)

	resp := req.Execute(handler)

	if resp.Code != http.StatusCreated {
		t.Errorf("Expected status code %d, got %d", http.StatusCreated, resp.Code)
	}

	// Verify the user was created in the database
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM users WHERE name = ? AND email = ?",
		"Test User", "test@example.com").Scan(&count)

	if err != nil {
		t.Fatalf("Failed to query users: %v", err)
	}

	if count != 1 {
		t.Errorf("Expected 1 user to be created, got %d", count)
	}
}

// Helper function to create a test handler
func createTestHandler(db *sql.DB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" && r.URL.Path == "/users" {
			// Parse the request body
			buf := make([]byte, r.ContentLength)
			r.Body.Read(buf)

			// In a real application, you would parse the JSON
			// For this test, we'll just insert a hardcoded user
			_, err := db.Exec("INSERT INTO users (name, email) VALUES (?, ?)",
				"Test User", "test@example.com")

			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			w.WriteHeader(http.StatusCreated)
			return
		}

		w.WriteHeader(http.StatusNotFound)
	})
}
