package handlers

import (
	"database/sql"
	"encoding/json"
	"go-inventory-management-api/config"
	"go-inventory-management-api/database"
	"log"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type Credentials struct {
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
}
type UserResponse struct {
	UserID   string `json:"id"`
	Username string `json:"username"`
	Token string `json:"token"`
}
type ErrorResponse struct {
	Message string `json:"message"`
}

type Claims struct{
	Username string `json:"username"`
	UserID string `json:"id"`
	jwt.RegisteredClaims
}

func respondWithError(w http.ResponseWriter, code int, message string){
	w.Header().Set("Content-Type", "application/json")
  w.WriteHeader(code)
  json.NewEncoder(w).Encode(ErrorResponse{Message:message})
}

func generateToken(username, userID string)(string,error){
	expirationTime := time.Now().Add(1*time.Hour)

	claims := &Claims{
		Username: username,
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),

		},
	}

	cfg, err := config.LoadConfig()
	if err != nil {
    log.Fatalf("Failed to load config: %v", err)
  }

	token :=jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(cfg.JwtSecret)
	if err != nil {
		return "", err
  }
	return tokenString, nil
}

// function to handle user registration
func Register(w http.ResponseWriter, r *http.Request) {
	var credentials Credentials
	err := json.NewDecoder(r.Body).Decode(&credentials)
	if err !=nil{
		// w.Header().Set("Content-Type", "application/json")
		// w.WriteHeader(http.StatusBadRequest)
		// json.NewEncoder(w).Encode(ErrorResponse{Message:"Invalid credentials"})
		respondWithError(w, http.StatusBadRequest, "Invalid credentials")
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(credentials.Password), bcrypt.DefaultCost)
	if err !=nil{
		respondWithError(w, http.StatusInternalServerError, "Internal server error")
		return
	}
	var userID string
	query :=`INSERT INTO "users" (username, password, is_admin, created_at, updated_at) 
	VALUES ($1, $2, false, NOW(), NOW()) RETURNING id`
	err = database.DB.QueryRow(query, credentials.Username, string(hashedPassword)).Scan(&userID)
	if err !=nil{
    respondWithError(w, http.StatusInternalServerError, "Error creating user")
    return
  }





	w.Header().Set("Content-Type", "application/json")
	// w.WriteHeader(http.StatusBadRequest)
	json.NewEncoder(w).Encode(UserResponse{UserID:userID, Username: credentials.Username})
	
}

//function to login user
func Login(w http.ResponseWriter, r *http.Request){
	var credentials Credentials
	err := json.NewDecoder(r.Body).Decode(&credentials)
	if err !=nil{
		// w.Header().Set("Content-Type", "application/json")
		// w.WriteHeader(http.StatusBadRequest)
		// json.NewEncoder(w).Encode(ErrorResponse{Message:"Invalid credentials"})
		respondWithError(w, http.StatusBadRequest, "Invalid credentials")
		return
	}

	var storedCredentials Credentials
	var userID string
	err = database.DB.QueryRow(`SELECT id, username, password FROM "users" WHERE username=$1`, 
				credentials.Username).Scan(&userID, &storedCredentials.Username, &storedCredentials.Password)
	
	if err != nil{
		if err == sql.ErrNoRows{
			respondWithError(w, http.StatusUnauthorized, "Invalid username or password")
			return
		}
		respondWithError(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(storedCredentials.Password), []byte(credentials.Password))
	if err != nil{
		respondWithError(w, http.StatusUnauthorized, "Invalid username or password")
    return
	}

	tokenString, err := generateToken(credentials.Username, userID)
	if err != nil{
		respondWithError(w, http.StatusInternalServerError, "Error generating token")
    return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(UserResponse{UserID: userID, Username: credentials.Username, Token: tokenString})
}