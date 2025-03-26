package handlers

import (
	"database/sql"
	"encoding/json"
	"go-inventory-management-api/database"
	"go-inventory-management-api/models"
	"go-inventory-management-api/utils"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func CreateItem(w http.ResponseWriter, r *http.Request) {
	var item models.InventoryItem
	err := json.NewDecoder(r.Body).Decode(&item)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	claims:= r.Context().Value("claims").(*utils.Claims)
	isAdmin := claims.IsAdmin

	if(!isAdmin){
		http.Error(w, "Not authorized to create", http.StatusUnauthorized)
		return
	}

	query := `INSERT INTO inventory_items (name, description, quantity, created_at, updated_at) 
			  VALUES ($1, $2, $3, NOW(), NOW()) RETURNING id`

	err = database.DB.QueryRow(query, item.Name, item.Description, item.Quantity).Scan(&item.ID)
	if err != nil {
		log.Printf("Failed to create inventory item: %v", err)
		http.Error(w, "Failed to create item", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(item)
}

func GetAllItems(w http.ResponseWriter, r *http.Request) {
	rows, err := database.DB.Query("SELECT id, name, description, quantity, created_at, updated_at FROM inventory_items")
	if err != nil {
		http.Error(w, "Failed to retrieve items", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var items []models.InventoryItem
	for rows.Next() {
		var item models.InventoryItem
		err := rows.Scan(&item.ID, &item.Name, &item.Description, &item.Quantity, &item.CreatedAt, &item.UpdatedAt)
		if err != nil {
			http.Error(w, "Failed to parse items", http.StatusInternalServerError)
			return
		}
		items = append(items, item)
	}

	json.NewEncoder(w).Encode(items)
}

func GetItem(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var item models.InventoryItem
	err := database.DB.QueryRow("SELECT id, name, description, quantity, created_at, updated_at FROM inventory_items WHERE id = $1", id).
		Scan(&item.ID, &item.Name, &item.Description, &item.Quantity, &item.CreatedAt, &item.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Item not found", http.StatusNotFound)
		} else {
			http.Error(w, "Failed to retrieve item", http.StatusInternalServerError)
		}
		return
	}

	json.NewEncoder(w).Encode(item)
}

func UpdateItem(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var item models.InventoryItem
	err := json.NewDecoder(r.Body).Decode(&item)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	query := `UPDATE inventory_items SET name = $1, description = $2, quantity = $3, updated_at = NOW() WHERE id = $4`
	_, err = database.DB.Exec(query, item.Name, item.Description, item.Quantity, id)
	if err != nil {
		http.Error(w, "Failed to update item", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Item updated successfully"})
}

func DeleteItem(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	query := `DELETE FROM inventory_items WHERE id = $1`
	_, err := database.DB.Exec(query, id)
	if err != nil {
		http.Error(w, "Failed to delete item", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Item deleted successfully"})
}
