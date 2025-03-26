package handlers

import (
	"encoding/json"
	"go-inventory-management-api/database"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

type Credentials struct {
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
}
type UserResponse struct {
	UserID   string `json:"id"`
	Username string `json:"username"`
}

type ErrorResponse struct {
	Message string `json:"message"`
}

func respondWithError(w http.ResponseWriter, code int, message string){
	w.Header().Set("Content-Type", "application/json")
  w.WriteHeader(code)
  json.NewEncoder(w).Encode(ErrorResponse{Message:message})
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