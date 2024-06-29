package main

import (
	"log"
	"net/http"

	"ip2country/config"
	"ip2country/handlers"
	"ip2country/ipdb"
	"ip2country/middleware"

	"github.com/gorilla/mux"
)

func main() {
	cfg := config.LoadConfig()

	db := &ipdb.CSVDatabase{}
	err := db.Load(cfg.IPDBPath)
	if err != nil {
		log.Fatalf("Could not load IP database: %s\n", err.Error())
	}

	router := mux.NewRouter()
	router.Use(middleware.RateLimiter(cfg.RateLimit))

	// Pass the database to the handler
	router.HandleFunc("/v1/find-country", handlers.MakeFindCountryHandler(db)).Methods("GET")

	log.Println("Starting server on port", cfg.Port)
	if err := http.ListenAndServe(":"+cfg.Port, router); err != nil {
		log.Fatalf("Could not start server: %s\n", err.Error())
	}
}
