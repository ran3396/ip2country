package main

import (
	"log"
	"net/http"

	"ip2country/config"
	"ip2country/ipdb"
	"ip2country/routes"
)

func main() {
	cfg := config.LoadConfig()

	db := &ipdb.CSVDatabase{}
	err := db.Load(cfg.IPDBPath)
	if err != nil {
		log.Fatalf("Could not load IP database: %s\n", err.Error())
	}

	// Setup the router using the routes package
	router := routes.SetupRoutes(db, cfg.RateLimit)

	log.Println("Starting server on port", cfg.Port)
	if err := http.ListenAndServe(":"+cfg.Port, router); err != nil {
		log.Fatalf("Could not start server: %s\n", err.Error())
	}
}
