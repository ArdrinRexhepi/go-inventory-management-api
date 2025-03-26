package routes

import (
	"go-inventory-management-api/handlers"
	"go-inventory-management-api/middleware"
	"log"
	"net/http"

	"go-inventory-management-api/utils"

	"github.com/gorilla/mux"
	"github.com/justinas/alice"
)

// type App struct{
// 	// DB *sql.DB
// 	JWTKey []byte
// }

func SetupRoutes(app *utils.App)  *mux.Router{
	router := mux.NewRouter()
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Println(w, "Inventory Management API is running!")
	}).Methods("GET")


	//middleware chain for user routes
	userChain := alice.New(middleware.AuthMiddleware(app), middleware.LoggingMiddleware)
	
	//middleware chain for admin routes
	adminChain := userChain.Append(middleware.AdminMiddleware)

	// Authentication routes
	router.Handle("/auth/register", alice.New(middleware.LoggingMiddleware).ThenFunc(handlers.Register)).Methods("POST")
	router.Handle("/auth/login", alice.New(middleware.LoggingMiddleware).ThenFunc(handlers.Login)).Methods("POST")

	//Inventory Items routes only token
	router.Handle("/inventory-items", userChain.ThenFunc(handlers.GetAllItems)).Methods("GET")
	router.Handle("/inventory-items/{id}", userChain.ThenFunc(handlers.GetItem)).Methods("GET")

	//Inventory routes that require admin access
	router.Handle("/inventory-items", adminChain.ThenFunc(handlers.CreateItem)).Methods("POST")
	router.Handle("/inventory-items", adminChain.ThenFunc(handlers.UpdateItem)).Methods("PUT")
	router.Handle("/inventory-items", adminChain.ThenFunc(handlers.DeleteItem)).Methods("DELETE")

	return router
}