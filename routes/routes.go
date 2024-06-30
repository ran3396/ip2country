package routes

import (
	"ip2country/handlers"
	"ip2country/ipdb"
	"ip2country/middleware"

	"github.com/gorilla/mux"
)

// SetupRoutes initializes the routes for the application
func SetupRoutes(db ipdb.IPDatabase, rateLimit int) *mux.Router {
	router := mux.NewRouter()

	// Apply the RateLimiter middleware to the router
	router.Use(middleware.RateLimiter(rateLimit))

	// Define the routes and handlers
	router.HandleFunc("/api/v1/find-country", handlers.MakeFindCountryHandler(db)).Methods("GET")

	return router
}
