package test

import (
	"os"
)

// Environment constants
const (
	EnvDevelopment = "development"
	EnvLocal       = "local"
	EnvProduction  = "production"
	EnvStaging     = "staging"
	EnvTesting     = "testing"
)

// TestConfig contains configuration options for setting up a test environment
type TestConfig struct {
	// Application settings
	AppName string
	AppURL  string
	AppEnv  string

	// Database settings
	DbDriver   string
	DbHost     string
	DbPort     string
	DbDatabase string
	DbUsername string
	DbPassword string

	// Server settings
	ServerHost string
	ServerPort string

	// Mail settings
	MailDriver   string
	MailHost     string
	MailPort     string
	MailUsername string
	MailPassword string
	EmailFrom    string
	EmailName    string

	// Security settings
	EnvEncryptionKey string
	VaultKey         string

	// Additional settings can be added as needed
	AdditionalEnvVars map[string]string
}

// DefaultTestConfig returns a default test configuration suitable for most test cases
func DefaultTestConfig() *TestConfig {
	return &TestConfig{
		AppName: "TEST APP",
		AppURL:  "http://localhost:8080",
		AppEnv:  EnvTesting,

		DbDriver:   "sqlite",
		DbHost:     "",
		DbPort:     "",
		DbDatabase: "file::memory:?cache=shared",
		DbUsername: "",
		DbPassword: "",

		ServerHost: "localhost",
		ServerPort: "8080",

		MailDriver:   "smtp",
		MailHost:     "127.0.0.1",
		MailPort:     "25",
		MailUsername: "",
		MailPassword: "",
		EmailFrom:    "test@example.com",
		EmailName:    "Test App",

		EnvEncryptionKey: "test_encryption_key_12345",
		VaultKey:         "abcdefghijklmnopqrstuvwxyz1234567890",

		AdditionalEnvVars: make(map[string]string),
	}
}

// SetupTestEnvironment configures the environment variables for testing based on the provided configuration
func SetupTestEnvironment(config *TestConfig) {
	// Application settings
	os.Setenv("APP_NAME", config.AppName)
	os.Setenv("APP_URL", config.AppURL)
	os.Setenv("APP_ENV", config.AppEnv)

	// Database settings
	os.Setenv("DB_DRIVER", config.DbDriver)
	os.Setenv("DB_HOST", config.DbHost)
	os.Setenv("DB_PORT", config.DbPort)
	os.Setenv("DB_DATABASE", config.DbDatabase)
	os.Setenv("DB_USERNAME", config.DbUsername)
	os.Setenv("DB_PASSWORD", config.DbPassword)

	// Server settings
	os.Setenv("SERVER_HOST", config.ServerHost)
	os.Setenv("SERVER_PORT", config.ServerPort)

	// Mail settings
	os.Setenv("MAIL_DRIVER", config.MailDriver)
	os.Setenv("MAIL_HOST", config.MailHost)
	os.Setenv("MAIL_PORT", config.MailPort)
	os.Setenv("MAIL_USERNAME", config.MailUsername)
	os.Setenv("MAIL_PASSWORD", config.MailPassword)
	os.Setenv("EMAIL_FROM_ADDRESS", config.EmailFrom)
	os.Setenv("EMAIL_FROM_NAME", config.EmailName)

	// Security settings
	os.Setenv("ENV_ENCRYPTION_KEY", config.EnvEncryptionKey)
	os.Setenv("VAULT_KEY", config.VaultKey)

	// Set any additional environment variables
	for key, value := range config.AdditionalEnvVars {
		os.Setenv(key, value)
	}
}

// CleanupTestEnvironment unsets all environment variables set by SetupTestEnvironment
func CleanupTestEnvironment(config *TestConfig) {
	// Application settings
	os.Unsetenv("APP_NAME")
	os.Unsetenv("APP_URL")
	os.Unsetenv("APP_ENV")

	// Database settings
	os.Unsetenv("DB_DRIVER")
	os.Unsetenv("DB_HOST")
	os.Unsetenv("DB_PORT")
	os.Unsetenv("DB_DATABASE")
	os.Unsetenv("DB_USERNAME")
	os.Unsetenv("DB_PASSWORD")

	// Server settings
	os.Unsetenv("SERVER_HOST")
	os.Unsetenv("SERVER_PORT")

	// Mail settings
	os.Unsetenv("MAIL_DRIVER")
	os.Unsetenv("MAIL_HOST")
	os.Unsetenv("MAIL_PORT")
	os.Unsetenv("MAIL_USERNAME")
	os.Unsetenv("MAIL_PASSWORD")
	os.Unsetenv("EMAIL_FROM_ADDRESS")
	os.Unsetenv("EMAIL_FROM_NAME")

	// Security settings
	os.Unsetenv("ENV_ENCRYPTION_KEY")
	os.Unsetenv("VAULT_KEY")

	// Unset any additional environment variables
	for key := range config.AdditionalEnvVars {
		os.Unsetenv(key)
	}
}
