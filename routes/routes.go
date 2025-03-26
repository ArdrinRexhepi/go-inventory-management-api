package routes

import (
	"go-inventory-management-api/handlers"
	"go-inventory-management-api/middleware"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/justinas/alice"
)

func SetupRoutes() *mux.Router{
	router := mux.NewRouter()
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Println(w, "Inventory Management API is running!")
	}).Methods("GET")

	router.Handle("/inventory-items", alice.New(middleware.LoggingMiddleware).ThenFunc(handlers.CreateItem)).Methods("POST")
	router.Handle("/inventory-items", alice.New(middleware.LoggingMiddleware).ThenFunc(handlers.GetAllItems)).Methods("GET")
	router.Handle("/inventory-items/{id}", alice.New(middleware.LoggingMiddleware).ThenFunc(handlers.GetItem)).Methods("GET")
	router.Handle("/inventory-items/{id}", alice.New(middleware.LoggingMiddleware).ThenFunc(handlers.UpdateItem)).Methods("PUT")
	router.Handle("/inventory-items/{id}", alice.New(middleware.LoggingMiddleware).ThenFunc(handlers.DeleteItem)).Methods("DELETE")

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Println(w, "Inventory Management API is running!")
	})

	return router
}