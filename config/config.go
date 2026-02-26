package config

import (
	"log"
	"os"
	"strconv"
)

// Config menyimpan semua konfigurasi aplikasi
type Config struct {
	HTTPPort       string
	GRPCPort       string
	DBDriver       string
	DBPath         string
	JWTSecret      string
	JWTExpiryHours int
	Environment    string
}

// LoadConfig membaca konfigurasi dari environment variables
func LoadConfig() *Config {
	return &Config{
		HTTPPort:       getEnv("HTTP_PORT", "8080"),
		GRPCPort:       getEnv("GRPC_PORT", "50051"),
		DBDriver:       getEnv("DB_DRIVER", "sqlite"),
		DBPath:         getEnv("DB_PATH", "test.db"),
		JWTSecret:      getEnv("JWT_SECRET", "default-secret-key-change-in-production"),
		JWTExpiryHours: getEnvAsInt("JWT_EXPIRY_HOURS", 24),
		Environment:    getEnv("ENV", "development"),
	}
}

// getEnv membaca environment variable dengan default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvAsInt membaca environment variable sebagai integer
func getEnvAsInt(key string, defaultValue int) int {
	valueStr := getEnv(key, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return defaultValue
}

// IsDevelopment mengecek apakah environment adalah development
func (c *Config) IsDevelopment() bool {
	return c.Environment == "development"
}

// IsProduction mengecek apakah environment adalah production
func (c *Config) IsProduction() bool {
	return c.Environment == "production"
}

// ValidateConfig memvalidasi konfigurasi yang diperlukan
func (c *Config) ValidateConfig() {
	if c.JWTSecret == "default-secret-key-change-in-production" && c.IsProduction() {
		log.Fatal("JWT_SECRET harus diset di production environment!")
	}
}
