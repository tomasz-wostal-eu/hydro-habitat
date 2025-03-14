package db

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq" // PostgreSQL driver
	"github.com/tomasz-wostal-eu/hydro-habitat/internal/config"
)

// NewPostgresConnection creates a new connection to the PostgreSQL database
func NewPostgresConnection(cfg config.DatabaseConfig) (*sql.DB, error) {
	// Construct the connection string
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName, cfg.SSLMode,
	)

	// Try to connect to the database with retries
	var db *sql.DB
	var err error

	for attempts := 1; attempts <= 10; attempts++ {
		db, err = sql.Open("postgres", dsn)
		if err == nil {
			// Check if connection is actually working
			err = db.Ping()
			if err == nil {
				break
			}
		}

		fmt.Printf("Attempt %d: Failed to connect to database: %v\n", attempts, err)
		time.Sleep(2 * time.Second) // Wait before retrying
	}

	if err != nil {
		return nil, fmt.Errorf("failed to connect to database after multiple attempts: %w", err)
	}

	// Set connection pool parameters
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	return db, nil
}
