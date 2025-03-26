package utils

import "github.com/golang-jwt/jwt/v5"

type App struct {
	// DB *sql.DB
	JWTKey []byte
}


type Claims struct{
	Username string `json:"username"`
	UserID string `json:"id"`
	IsAdmin bool `json:"is_admin"`
	jwt.RegisteredClaims
}