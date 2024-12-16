package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

func LoadEnv() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found. Loading environment variables from OS.")
		os.Exit(1)
	}
}

// GetEnv retrieves the value of an environment variable by name.
// Returns the default value if the variable is not set.
func GetEnv(key string, defaultValue string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}
	return value
}

// GetEnvAsInt retrieves an environment variable as an integer.
// Returns the default value if the variable is not set or cannot be converted.
func GetEnvAsInt(key string, defaultValue int) int {
	valueStr, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}
	value, err := strconv.Atoi(valueStr)
	if err != nil {
		log.Printf("Error converting %s to int: %v. Using default value: %d\n", key, err, defaultValue)
		return defaultValue
	}
	return value
}

// GetEnvAsBool retrieves an environment variable as a boolean.
// Returns the default value if the variable is not set or cannot be converted.
func GetEnvAsBool(key string, defaultValue bool) bool {
	valueStr, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}
	value, err := strconv.ParseBool(valueStr)
	if err != nil {
		log.Printf("Error converting %s to bool: %v. Using default value: %v\n", key, err, defaultValue)
		return defaultValue
	}
	return value
}
