# Testing Utilities

This package provides a set of utilities for testing Go applications in the Dracory ecosystem. It includes tools for setting up test environments, managing test databases, and testing HTTP endpoints.

## Key Components

### Test Configuration

The `test_config.go` file provides utilities for setting up test environments:

- `TestConfig`: A struct that contains configuration options for setting up a test environment
- `DefaultTestConfig()`: Returns a default test configuration suitable for most test cases
- `SetupTestEnvironment()`: Configures environment variables for testing
- `CleanupTestEnvironment()`: Cleans up environment variables after testing

### Test Database

The `test_db.go` file provides utilities for setting up and managing test databases:

- `NewTestDB()`: Creates a new test database connection (defaults to in-memory SQLite)
- `CloseTestDB()`: Safely closes a test database connection
- `ExecuteSQL()`: Executes SQL statements on the database
- `CreateTestTable()`: Creates test tables in the database
- `DropTestTable()`: Drops test tables from the database

### Test HTTP

The `test_http.go` file provides utilities for testing HTTP endpoints:

- `TestHTTPRequest`: A struct for building test HTTP requests
- `TestHTTPServer`: A wrapper around httptest.Server for testing HTTP servers
- Helper methods for executing HTTP requests and handling responses

### Test Key

The `test_key.go` file provides a utility for generating test keys:

- `TestKey()`: Generates a consistent hash based on database configuration

## Usage Examples

### Setting Up a Test Environment

```go
import (
    "testing"
    "github.com/dracory/base/test"
)

func TestSomething(t *testing.T) {
    // Create and customize a test configuration
    config := testutils.DefaultTestConfig()
    config.AppName = "My Test App"

    // Set up the test environment
    testutils.SetupTestEnvironment(config)

    // Run your tests...

    // Clean up the test environment
    testutils.CleanupTestEnvironment(config)
}
```

### Using the Test Database

```go
import (
    "testing"
    "github.com/dracory/base/test"
)

func TestWithDatabase(t *testing.T) {
    // Create a test database
    db, err := testutils.NewTestDB(nil) // Uses default in-memory SQLite
    if err != nil {
        t.Fatalf("Failed to create test database: %v", err)
    }
    defer testutils.CloseTestDB(db)

    // Create a test table
    err = testutils.CreateTestTable(db, "users", "id INTEGER PRIMARY KEY, name TEXT, email TEXT")
    if err != nil {
        t.Fatalf("Failed to create test table: %v", err)
    }

    // Execute SQL
    err = testutils.ExecuteSQLWithArgs(db, "INSERT INTO users (name, email) VALUES (?, ?)", "Test User", "test@example.com")
    if err != nil {
        t.Fatalf("Failed to insert test data: %v", err)
    }

    // Run your tests with the database...
}
```

### Testing HTTP Endpoints

```go
import (
    "net/http"
    "testing"
    "github.com/dracory/base/test"
)

func TestHTTPEndpoint(t *testing.T) {
    // Create a test handler
    handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusOK)
        w.Write([]byte("Hello, world!"))
    })

    // Create a request
    req := test.NewTestHTTPRequest("GET", "/hello")
        .WithHeader("X-Test", "test-value")

    // Execute the request
    resp := req.Execute(handler)

    // Check the response
    if resp.Code != http.StatusOK {
        t.Errorf("Expected status code %d, got %d", http.StatusOK, resp.Code)
    }

    if resp.Body.String() != "Hello, world!" {
        t.Errorf("Expected body %q, got %q", "Hello, world!", resp.Body.String())
    }
}
```

## Best Practices

1. Always clean up resources after tests (close database connections, clean up environment variables)
2. Use in-memory SQLite for tests to avoid external dependencies
3. Use real database connections and stores instead of mocks
4. Keep tests isolated and independent
5. Use the provided utilities to simplify test setup and teardown
