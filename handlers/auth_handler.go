package handlers

import (
	"database/sql"
	"encoding/json"
	"go-inventory-management-api/internal/config"
	"go-inventory-management-api/internal/database"
	"go-inventory-management-api/utils"
	"log"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/mux"
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

func generateToken(username string, userID string, isAdmin bool)(string,error){
	expirationTime := time.Now().Add(1*time.Hour)

	claims := &utils.Claims{
		Username: username,
		UserID: userID,
		IsAdmin: isAdmin,
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
		http.Error(w, "Invalid credentials", http.StatusBadRequest)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(credentials.Password), bcrypt.DefaultCost)
	if err !=nil{
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	var userID string
	query :=`INSERT INTO "users" (username, password, is_admin, created_at, updated_at) 
	VALUES ($1, $2, false, NOW(), NOW()) RETURNING id`
	err = database.DB.QueryRow(query, credentials.Username, string(hashedPassword)).Scan(&userID)
	if err !=nil{
		http.Error(w, "Error creating user", http.StatusInternalServerError)
    return
  }
	tokenString, err := generateToken(credentials.Username, userID, false)
	if err != nil{
		http.Error(w, "Error generating token", http.StatusInternalServerError)
    return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(UserResponse{UserID: userID, Username: credentials.Username, Token: tokenString})
	
}

//function to login user
func Login(w http.ResponseWriter, r *http.Request){
	var credentials Credentials
	err := json.NewDecoder(r.Body).Decode(&credentials)
	if err !=nil{
		http.Error(w, "Invalid credentials", http.StatusBadRequest)
		return
	}
	var storedCredentials Credentials
	var userID string
	var isAdmin bool
	err = database.DB.QueryRow(`SELECT id, username, is_admin, password FROM "users" WHERE username=$1`, 
				credentials.Username).Scan(&userID, &storedCredentials.Username, &isAdmin, &storedCredentials.Password)
	
	if err != nil{
		if err == sql.ErrNoRows{
			http.Error(w, "Invalid username or password", http.StatusUnauthorized)
			return
		}
		http.Error(w, "Failed to parse items", http.StatusInternalServerError)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(storedCredentials.Password), []byte(credentials.Password))
	if err != nil{
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
    return
	}

	tokenString, err := generateToken(credentials.Username, userID, isAdmin)
	if err != nil{
		http.Error(w, "Error generating token", http.StatusInternalServerError)
    return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(UserResponse{UserID: userID, Username: credentials.Username, Token: tokenString})
}

//function to turn a user into an admin
func MakeNewAdmin(w http.ResponseWriter, r *http.Request){
	vars := mux.Vars(r)
	userID := vars["id"]

	var userExists bool
	err:= database.DB.QueryRow(`SELECT EXISTS(SELECT 1 FROM users WHERE id=$1)`, userID).Scan(&userExists)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	if !userExists {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}
	var isAdmin bool
	err = database.DB.QueryRow(`SELECT is_admin FROM users WHERE id = $1`, userID).Scan(&isAdmin)
	if err != nil {
		http.Error(w, "Failed to check user admin status", http.StatusInternalServerError)
		return
	}
	if isAdmin {
		http.Error(w, "User is already an admin", http.StatusConflict)
    return
	}

	_, err = database.DB.Exec(`UPDATE users SET is_admin=true where id=$1`,userID)
	if err != nil{
		http.Error(w, "Failed to make user admin", http.StatusInternalServerError)
    return
	}

	w.WriteHeader(http.StatusNoContent)
}