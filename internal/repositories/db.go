package repositories

import (
	"fmt"
	"log"
	"os"

	"database/sql"
	_ "github.com/lib/pq"
)

// InitAndCloseDB initializes the database and ensures it gets closed on exit.
func InitAndCloseDB() (*sql.DB, func(), error) {
	databaseURL := os.Getenv("DATABASE_URL")
	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to connect to the database: %v", err)
	}

	// Verify the connection
	if err := db.Ping(); err != nil {
		db.Close()
		return nil, nil, fmt.Errorf("failed to ping the database: %v", err)
	}

	// Close db function to defer.
	cleanup := func() {
		if err := db.Close(); err != nil {
			log.Printf("error closing database connection: %v", err)
		}
	}

	return db, cleanup, nil
}
