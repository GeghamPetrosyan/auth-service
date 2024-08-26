package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
	"time"
)

type Config struct {
	DBHost          string
	DBPort          string
	DBUser          string
	DBPassword      string
	DBName          string
	JWTSecret       string
	AccessTokenTTL  time.Duration
	RefreshTokenTTL time.Duration
	EmailService    string
}

// LoadConfig загружает переменные окружения из .env файла
func LoadConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file, using environment variables")
	}

	config := &Config{
		DBHost:          getEnv("DB_HOST", "localhost"),
		DBPort:          getEnv("DB_PORT", "5432"),
		DBUser:          getEnv("DB_USER", "postgres"),
		DBPassword:      getEnv("DB_PASSWORD", "password"),
		DBName:          getEnv("DB_NAME", "authdb"),
		JWTSecret:       getEnv("JWT_SECRET", "your-secret-key"),
		AccessTokenTTL:  time.Hour * 1,      // Время жизни Access токена
		RefreshTokenTTL: time.Hour * 24 * 7, // Время жизни Refresh токена
		EmailService:    getEnv("EMAIL_SERVICE", "mock"),
	}

	return config
}

// getEnv читает переменную окружения или возвращает значение по умолчанию
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
