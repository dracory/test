package test

import (
	"database/sql"
	"fmt"

	_ "github.com/glebarez/sqlite"
)

// DBConfig contains configuration for test database
type DBConfig struct {
	Driver   string
	Host     string
	Port     string
	Database string
	Username string
	Password string
}

// DefaultDBConfig returns a default SQLite in-memory database configuration
func DefaultDBConfig() *DBConfig {
	return &DBConfig{
		Driver:   "sqlite",
		Database: "file::memory:?cache=shared",
	}
}

// NewTestDB creates a new test database connection
// By default, it creates an in-memory SQLite database
func NewTestDB(config *DBConfig) (*sql.DB, error) {
	if config == nil {
		config = DefaultDBConfig()
	}

	var dsn string
	switch config.Driver {
	case "sqlite":
		dsn = config.Database
	case "mysql":
		dsn = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
			config.Username, config.Password, config.Host, config.Port, config.Database)
	case "postgres":
		dsn = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
			config.Host, config.Port, config.Username, config.Password, config.Database)
	default:
		return nil, fmt.Errorf("unsupported database driver: %s", config.Driver)
	}

	db, err := sql.Open(config.Driver, dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	// Verify the connection
	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return db, nil
}

// CloseTestDB safely closes a test database connection
func CloseTestDB(db *sql.DB) error {
	if db != nil {
		return db.Close()
	}
	return nil
}

// ExecuteSQL executes a SQL statement on the database
func ExecuteSQL(db *sql.DB, sql string) error {
	_, err := db.Exec(sql)
	return err
}

// ExecuteSQLWithArgs executes a SQL statement with arguments on the database
func ExecuteSQLWithArgs(db *sql.DB, sql string, args ...interface{}) error {
	_, err := db.Exec(sql, args...)
	return err
}

// CreateTestTable creates a test table in the database
func CreateTestTable(db *sql.DB, tableName string, schema string) error {
	createTableSQL := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (%s)", tableName, schema)
	return ExecuteSQL(db, createTableSQL)
}

// DropTestTable drops a test table from the database
func DropTestTable(db *sql.DB, tableName string) error {
	dropTableSQL := fmt.Sprintf("DROP TABLE IF EXISTS %s", tableName)
	return ExecuteSQL(db, dropTableSQL)
}
