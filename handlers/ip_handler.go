package handlers

import (
	"log"
	"net/http"

	"ip2country/ipdb"
	"ip2country/utils"
)

// Return a handler function that uses the given database to find the country for the given IP.
// If the IP is not found, respond with a 404 status code.
// If the IP is missing, respond with a 400 status code.
// If there is an error querying the database, respond with a 500 status code.
// If the IP is found, respond with a 200 status code and a JSON object containing the IP, city, and country.
func MakeFindCountryHandler(db ipdb.IPDatabase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ip := r.URL.Query().Get("ip")

		log.Println("Received request to find", ip)

		if ip == "" {
			utils.RespondWithError(w, http.StatusBadRequest, "IP parameter is missing")
			return
		}

		city, country, err := db.Find(ip)
		if err == ipdb.ErrNotFound {
			message := "IP not found in the database"
			log.Println(message)
			utils.RespondWithError(w, http.StatusNotFound, message)
			return
		} else if err != nil {
			message := "Failed to query IP database"
			log.Println(message)
			utils.RespondWithError(w, http.StatusInternalServerError, message)
			return
		}

		ipInfo := map[string]string{
			"ip":      ip,
			"city":    city,
			"country": country,
		}

		log.Println("Found mapping for", ip, ":", city, ",", country)

		utils.RespondWithJSON(w, http.StatusOK, ipInfo)
	}
}
