package test

import (
	"github.com/dracory/str"
)

// TestKey is a pseudo secret test key used for testing specific unit cases
// where a secret key is required but not available in the testing environment.
// It generates a consistent hash based on the provided database configuration.
func TestKey(dbDriver, dbHost, dbPort, dbName, dbUser, dbPass string) string {
	return str.MD5(dbDriver + dbHost + dbPort + dbName + dbUser + dbPass)
}
