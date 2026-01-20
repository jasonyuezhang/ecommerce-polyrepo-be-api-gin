package config

import (
	"os"
	"strconv"
)

// Config holds all configuration for the application
type Config struct {
	// Server settings
	Port        string
	Environment string

	// JWT settings
	JWTSecret     string
	JWTExpiration int // in hours

	// gRPC service addresses
	UserServiceAddr      string
	ListingServiceAddr   string
	InventoryServiceAddr string

	// CORS settings
	AllowedOrigins []string

	// Rate limiting
	RateLimit int // requests per second
}

// Load reads configuration from environment variables
func Load() *Config {
	return &Config{
		Port:                 getEnv("PORT", "8080"),
		Environment:          getEnv("ENVIRONMENT", "development"),
		JWTSecret:            getEnv("JWT_SECRET", "your-secret-key-change-in-production"),
		JWTExpiration:        getEnvAsInt("JWT_EXPIRATION_HOURS", 24),
		UserServiceAddr:      getEnv("USER_SERVICE_ADDR", "localhost:50051"),
		ListingServiceAddr:   getEnv("LISTING_SERVICE_ADDR", "localhost:50052"),
		InventoryServiceAddr: getEnv("INVENTORY_SERVICE_ADDR", "localhost:50053"),
		AllowedOrigins:       getEnvAsSlice("ALLOWED_ORIGINS", []string{"http://localhost:3000"}),
		RateLimit:            getEnvAsInt("RATE_LIMIT", 100),
	}
}

// getEnv gets an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

// getEnvAsInt gets an environment variable as an integer or returns a default value
func getEnvAsInt(key string, defaultValue int) int {
	if value, exists := os.LookupEnv(key); exists {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

// getEnvAsSlice gets an environment variable as a slice or returns a default value
func getEnvAsSlice(key string, defaultValue []string) []string {
	if value, exists := os.LookupEnv(key); exists && value != "" {
		// Simple comma-separated parsing
		var result []string
		start := 0
		for i := 0; i <= len(value); i++ {
			if i == len(value) || value[i] == ',' {
				if i > start {
					result = append(result, value[start:i])
				}
				start = i + 1
			}
		}
		if len(result) > 0 {
			return result
		}
	}
	return defaultValue
}
