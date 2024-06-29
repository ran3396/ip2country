package ipdb

import (
	"encoding/csv"
	"os"
	"sync"
)

// IPDatabase is the interface that different IP databases must implement.
type IPDatabase interface {
	Load(filePath string) error
	Find(ip string) (string, string, error)
}

// CSVDatabase is a simple implementation of IPDatabase using CSV files.
type CSVDatabase struct {
	mu       sync.RWMutex
	data     map[string][]string
	filePath string
}

// Load reads the CSV file and loads the data into memory.
func (db *CSVDatabase) Load(filePath string) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return err
	}

	db.data = make(map[string][]string)
	for _, record := range records {
		if len(record) >= 3 {
			ip := record[0]
			city := record[1]
			country := record[2]
			db.data[ip] = []string{city, country}
		}
	}
	db.filePath = filePath
	return nil
}

// Find searches for the given IP in the database and returns the city and country.
func (db *CSVDatabase) Find(ip string) (string, string, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	if location, exists := db.data[ip]; exists {
		return location[0], location[1], nil
	}
	return "", "", ErrNotFound
}

// ErrNotFound is returned when the IP is not found in the database.
var ErrNotFound = os.ErrNotExist
