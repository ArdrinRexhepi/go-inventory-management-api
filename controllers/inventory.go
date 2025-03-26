package controllers

import (
	"encoding/json"
	"go-inventory-management-api/database"
	"go-inventory-management-api/models"
	"net/http"
)

func CreateItem(w http.ResponseWriter, r *http.Request) {
	var item models.InventoryItem
	err := json.NewDecoder(r.Body).Decode(&item)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	query := `INSERT INTO inventory_items (name, description, quantity, created_at, updated_at, last_restock) 
			  VALUES ($1, $2, $3, NOW(), NOW(), NOW()) RETURNING id`

	err = database.DB.QueryRow(query, item.Name, item.Description, item.Quantity).Scan(&item.ID)
	if err != nil {
		http.Error(w, "Failed to create item", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(item)
}