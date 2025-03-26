package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

var DB *sql.DB

func ConnectDb(connStr string) {
	var err error
	DB, err = sql.Open("postgres", connStr)  // Use global DB variable
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	if err = DB.Ping(); err != nil {
		log.Fatalf("PING ERROR: %v", err)
	}

	fmt.Println("Successfully connected to database")
	runMigrations(DB)
}

func runMigrations(db *sql.DB) {
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		log.Fatalf("Migration driver error: %v", err)
	}

	m, err := migrate.NewWithDatabaseInstance("file://migrations", "postgres", driver)
	if err != nil {
		log.Fatalf("Migration error: %v", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatalf("Migration execution error: %v", err)
	}

	fmt.Println("Database migrations applied successfully")
}