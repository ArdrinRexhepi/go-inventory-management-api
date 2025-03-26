package config

import (
	"os"

	"github.com/lpernett/godotenv"
)

type Config struct {
	DatabaseURL string
	JwtSecret  []byte
}

func LoadConfig() (*Config, error) {
	err:= godotenv.Load()
	if err !=nil{
		return nil, err
	}

	JwtSecret := []byte(os.Getenv("JWT_SECRET"))
	if len(JwtSecret) == 0 {
    return nil, os.ErrInvalid
  }
	DatabaseURL := os.Getenv("DATABASE_URL")
	if len(DatabaseURL) == 0 {
    return nil, os.ErrInvalid
  }
	return &Config{
		DatabaseURL: os.Getenv("DATABASE_URL"),
		JwtSecret:   JwtSecret,
	}, nil
}