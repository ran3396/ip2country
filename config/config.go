package config

import (
	"os"
	"strconv"
)

type Config struct {
	Port      string
	RateLimit int
	IPDBPath  string
}

func LoadConfig() Config {
	return Config{
		Port:      getEnv("PORT", "8080"),
		RateLimit: getEnvAsInt("RATE_LIMIT", 10),
		IPDBPath:  getEnv("IP_DB_PATH", "testdata/ipdb.csv"),
	}
}

// getEnv is a helper function to get an environment variable with a default value.
func getEnv(key, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultVal
}

// getEnvAsInt is a helper function to get an integer environment variable.
func getEnvAsInt(name string, defaultVal int) int {
	if valueStr, exists := os.LookupEnv(name); exists {
		if value, err := strconv.Atoi(valueStr); err == nil {
			return value
		}
	}
	return defaultVal
}
