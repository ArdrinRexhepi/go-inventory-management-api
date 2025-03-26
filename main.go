package main

import (
	"go-inventory-management-api/config"
	"go-inventory-management-api/database"
	"go-inventory-management-api/routes"
	"log"
	"net/http"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
    log.Fatalf("Failed to load config: %v", err)
  }
	if err := database.ConnectDb(cfg.DatabaseURL); err != nil {
		log.Fatalf("Error initializing the database: %v", err)
		return
	}

	router := routes.SetupRoutes()

	log.Println("Listening on port", cfg.ServerPort)
	log.Fatal(http.ListenAndServe(":"+cfg.ServerPort, router))
}