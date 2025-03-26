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
	database.ConnectDb(cfg.DatabaseURL)

	router := routes.SetupRoutes()

	log.Println("Listening on port", cfg.ServerPort)
	log.Fatal(http.ListenAndServe(":"+cfg.ServerPort, router))
}