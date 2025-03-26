package routes

import (
	"go-inventory-management-api/handlers"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func SetupRoutes() *mux.Router{
	router := mux.NewRouter()
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Println(w, "Inventory Management API is running!")
	}).Methods("GET")

	router.HandleFunc("/inventory-items", handlers.CreateItem).Methods("POST")
	router.HandleFunc("/inventory-items", handlers.GetAllItems).Methods("GET")
	router.HandleFunc("/inventory-items/{id}", handlers.GetItem).Methods("GET")
	router.HandleFunc("/inventory-items/{id}", handlers.UpdateItem).Methods("PUT")
	router.HandleFunc("/inventory-items/{id}", handlers.DeleteItem).Methods("DELETE")

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Println(w, "Inventory Management API is running!")
	})

	return router
}