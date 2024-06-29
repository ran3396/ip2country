package handlers

import (
	"encoding/json"
	"ip2country/ipdb"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"testing"
	"time"
)

// setupDatabase is a helper function to load the IP database for testing.
func setupDatabase(t *testing.T) *ipdb.CSVDatabase {
	db := &ipdb.CSVDatabase{}
	os.Setenv("IP_DB_PATH", "../testdata/ipdb.csv")
	err := db.Load(os.Getenv("IP_DB_PATH"))
	if err != nil {
		t.Fatalf("Failed to load database: %v", err)
	}
	return db
}

// TestFindCountryValidIP tests the handler with a valid IP address.
// It expects a 200 status code and a JSON response with the IP, city, and country.
func TestFindCountryValidIP(t *testing.T) {
	db := setupDatabase(t)

	req, err := http.NewRequest("GET", "/v1/find-country?ip=2.22.233.255", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := MakeFindCountryHandler(db)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	expected := map[string]string{
		"ip":      "2.22.233.255",
		"city":    "London",
		"country": "United Kingdom",
	}
	var received map[string]string
	errJson := json.NewDecoder(rr.Body).Decode(&received)
	if errJson != nil {
		t.Fatalf("Failed to decode response body: %v", err)
	}
	if !reflect.DeepEqual(expected, received) {
		t.Errorf("handler returned unexpected body: got %v want %v",
			received, expected)
	}
}

// TestFindCountryMissingIP tests the handler with a missing IP address.
// It expects a 400 status code and an error message in the response.
func TestFindCountryMissingIP(t *testing.T) {
	db := setupDatabase(t)

	req, err := http.NewRequest("GET", "/v1/find-country", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := MakeFindCountryHandler(db)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusBadRequest)
	}

	expected := `{"error":"IP parameter is missing"}`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

// TestFindCountryInvalidIP tests the handler with an invalid IP address.
// It expects a 404 status code and an error message in the response.
func TestFindCountryInvalidIP(t *testing.T) {
	db := setupDatabase(t)

	req, err := http.NewRequest("GET", "/v1/find-country?ip=invalid_ip", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := MakeFindCountryHandler(db)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusNotFound)
	}

	expected := `{"error":"IP not found in the database"}`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

// TestRateLimiter tests the rate limiter middleware functionality.
// It expects the first request to return a 200 status code and the second request to return a 429 status code.
func TestRateLimiter(t *testing.T) {
	db := setupDatabase(t)
	os.Setenv("RATE_LIMIT", "1") // One request per second

	handler := MakeFindCountryHandler(db)
	testCases := []struct {
		IP       string
		Expected int
	}{
		{"2.22.233.255", http.StatusOK},
		{"2.22.233.255", http.StatusTooManyRequests},
	}

	for _, tc := range testCases {
		req, err := http.NewRequest("GET", "/v1/find-country?ip="+tc.IP, nil)
		if err != nil {
			t.Fatal(err)
		}
		rr := httptest.NewRecorder()

		// Execute
		handler.ServeHTTP(rr, req)

		// Validate
		if status := rr.Code; status != tc.Expected {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, tc.Expected)
		}

		// Add a delay to avoid immediate rate limit hit for next request
		time.Sleep(1 * time.Second)
	}
}
