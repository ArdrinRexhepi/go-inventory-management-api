package config

import (
	"os"

	"github.com/lpernett/godotenv"
)

type Config struct {
	DatabaseURL string
	ServerPort  string
	JwtSecret  string
}

func LoadConfig() (*Config, error) {
	err:= godotenv.Load()
	if err !=nil{
		return nil, err
	}
	return &Config{
		DatabaseURL: os.Getenv("DATABASE_URL"),
		JwtSecret: os.Getenv("JWT_SECRET"),
		ServerPort:  getEnv("SERVER_PORT", "5000"),
	}, nil
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}