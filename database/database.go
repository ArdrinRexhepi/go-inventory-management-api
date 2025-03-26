package database

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

var DB *sql.DB

func ConnectDb(connStr string) error {
	var err error
	DB, err = sql.Open("postgres", connStr)  // Use global DB variable
	if err != nil {
		return fmt.Errorf("failed to initialize database: %v", err)
	}

	if err = DB.Ping(); err != nil {
		return fmt.Errorf("ping error: %v", err)
	}

	fmt.Println("Successfully connected to database")
	runMigrations(DB)
	return nil
}

func runMigrations(db *sql.DB) error {
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("migration driver error: %v", err)
	}

	m, err := migrate.NewWithDatabaseInstance("file://migrations", "postgres", driver)
	if err != nil {
		return fmt.Errorf("migration error: %v", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("migration execution error: %v", err)
	}

	fmt.Println("Database migrations applied successfully")
	return nil
}