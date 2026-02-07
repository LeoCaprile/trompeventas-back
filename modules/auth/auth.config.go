package auth

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	JWTSecret            string
	FrontendURL          string
	BackendURL           string
	GoogleClientID       string
	GoogleClientSecret   string
	GoogleRedirectURL    string
	AccessTokenDuration  int // minutes
	RefreshTokenDuration int // days
}

var AppConfig *Config

func LoadConfig() {
	godotenv.Load()

	AppConfig = &Config{
		JWTSecret:            os.Getenv("JWT_SECRET"),
		FrontendURL:          getEnvOrDefault("FRONTEND_URL", "http://localhost:5173"),
		BackendURL:           getEnvOrDefault("BACKEND_URL", "http://localhost:8080"),
		GoogleClientID:       os.Getenv("GOOGLE_CLIENT_ID"),
		GoogleClientSecret:   os.Getenv("GOOGLE_CLIENT_SECRET"),
		AccessTokenDuration:  15,
		RefreshTokenDuration: 7,
	}

	AppConfig.GoogleRedirectURL = AppConfig.BackendURL + "/auth/oauth/google/callback"

	if AppConfig.JWTSecret == "" {
		log.Fatal("JWT_SECRET environment variable is required")
	}

	if AppConfig.GoogleClientID == "" {
		log.Println("WARNING: GOOGLE_CLIENT_ID not set, Google OAuth will not work")
	}
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
