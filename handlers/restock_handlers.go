package handlers

import (
	"encoding/json"
	"go-inventory-management-api/database"
	"go-inventory-management-api/models"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func RestockItem(w http.ResponseWriter, r *http.Request) {
	var restock models.RestockHistory
	
	err := json.NewDecoder(r.Body).Decode(&restock)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	itemID := restock.ItemID

	// Validating the restock amount
	if restock.RestockAmount < 10 || restock.RestockAmount > 1000 {
		http.Error(w, "Restock amount must be between 10 and 1000", http.StatusBadRequest)
		return
	}
	log.Println("444444444")

	// Check how much is the restock count in the last 24 hours
	var count int
	err = database.DB.QueryRow(
		`SELECT COUNT(*) FROM restock_history WHERE item_id=$1 AND restocked_at >= NOW() - INTERVAL '24 hours'`, 
		itemID,
	).Scan(&count)

	if err != nil {
		http.Error(w, "Failed to retrieve restock history", http.StatusInternalServerError)
		return
	}

	// Enforce rate limit (max 3 restocks in 24 hours)
	if count >= 3 {
		http.Error(w, "Too many restocks in the last 24 hours for this item", http.StatusTooManyRequests)
		return
	}

	//Added transaction since we are adding to restock_history and also updating with that quantity the inventory_item
	transaction, err := database.DB.Begin()
	if err != nil {
		http.Error(w, "Failed to start transaction", http.StatusInternalServerError)
		return
	}

	// Inserting the restock record data
	_, err = transaction.Exec(
		`INSERT INTO restock_history (item_id, restock_amount, restocked_at) VALUES ($1, $2, NOW())`,
		itemID, restock.RestockAmount,
	)
	if err != nil {
		transaction.Rollback()
		http.Error(w, "Failed to save restock record", http.StatusInternalServerError)
		return
	}

	//Updating the inventory_item quantity
	_, err = transaction.Exec(
    `UPDATE inventory_items SET quantity = quantity + $1 WHERE id = $2`,
    restock.RestockAmount, itemID,
  )
  if err != nil {
    transaction.Rollback()
    http.Error(w, "Failed to update inventory item quantity", http.StatusInternalServerError)
    return
  }

  err = transaction.Commit()
  if err != nil {
    http.Error(w, "Failed to commit transaction", http.StatusInternalServerError)
    return
  }

  restock.ID = itemID
  w.WriteHeader(http.StatusCreated)
  json.NewEncoder(w).Encode(restock)
}

func GetRestockHistory(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	rows, err := database.DB.Query(`SELECT * FROM restock_history where item_id=$1`,id)
	if err != nil {
		http.Error(w, "Failed to retrieve items", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var items []models.RestockHistory
	for rows.Next(){
		var item models.RestockHistory
		err := rows.Scan(&item.ID,&item.ItemID, &item.RestockAmount, &item.RestockedAt)
		if err != nil {
			http.Error(w, "Failed to parse items", http.StatusInternalServerError)
			return
		}
		items = append(items, item)
	}

	json.NewEncoder(w).Encode(items)
}