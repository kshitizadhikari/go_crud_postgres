package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Config holds all application configuration
type Config struct {
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	DBSSLMode  string
	AppPort    string
}

// LoadConfig loads configuration from .env file
func LoadConfig() *Config {
	// Load .env file (ignore error if file doesn't exist)
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, using environment variables")
	}

	return &Config{
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBUser:     getEnv("DB_USER", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", "password"),
		DBName:     getEnv("DB_NAME", "go_crud"),
		DBSSLMode:  getEnv("DB_SSLMODE", "disable"),
		AppPort:    getEnv("APP_PORT", "8080"),
	}
}

// GetDBURL returns the PostgreSQL connection string
// func (c *Config) GetDBURL() string {
// 	return "postgres://" + c.DBUser + ":" + c.DBPassword + "@" +
// 		c.DBHost + ":" + c.DBPort + "/" + c.DBName +
// 		"?sslmode=" + c.DBSSLMode
// }

// GetDSN returns DSN format for GORM
func (c *Config) GetDSN() string {
	return fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=Asia/Kathmandu",
		c.DBHost,
		c.DBUser,
		c.DBPassword,
		c.DBName,
		c.DBPort,
		c.DBSSLMode,
	)
}

// Helper function to get env variable with default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
