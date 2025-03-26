package main

import (
	"go-inventory-management-api/config"
	"go-inventory-management-api/database"
	"go-inventory-management-api/routes"
	"go-inventory-management-api/utils"
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
	defer database.DB.Close()

	app:=&utils.App{JWTKey:cfg.JwtSecret}

	router := routes.SetupRoutes(app)

	log.Println("Listening on port 5000")
	log.Fatal(http.ListenAndServe(":5000", router))
}