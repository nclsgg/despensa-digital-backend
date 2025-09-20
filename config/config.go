package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseURL   string
	JWTSecret     string
	JWTExpiration string
	JWTIssuer     string
	JWTAudience   string
	RedisURL      string
	RedisPassword string
	RedisDB       string
	RedisUsername string
	Port          string
	CorsOrigin    string
	
	// OAuth Config
	GoogleClientID     string
	GoogleClientSecret string
	GoogleCallbackURL  string
	FrontendURL        string
	SessionSecret      string
}

func LoadConfig() *Config {
	_ = godotenv.Load()

	cfg := &Config{
		DatabaseURL:   os.Getenv("DATABASE_URL"),
		JWTSecret:     os.Getenv("JWT_SECRET"),
		JWTExpiration: os.Getenv("JWT_EXPIRATION"),
		JWTIssuer:     os.Getenv("JWT_ISSUER"),
		JWTAudience:   os.Getenv("JWT_AUDIENCE"),
		RedisURL:      os.Getenv("REDIS_URL"),
		RedisPassword: os.Getenv("REDIS_PASSWORD"),
		RedisDB:       os.Getenv("REDIS_DB"),
		RedisUsername: os.Getenv("REDIS_USERNAME"),
		Port:          getEnv("PORT", "3030"),
		CorsOrigin:    getEnv("CORS_ORIGIN", "http://localhost:3000"),
		
		// OAuth Config
		GoogleClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		GoogleClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		GoogleCallbackURL:  getEnv("GOOGLE_CALLBACK_URL", "http://localhost:3030/auth/oauth/google/callback"),
		FrontendURL:        getEnv("FRONTEND_URL", "http://localhost:3000"),
		SessionSecret:      getEnv("SESSION_SECRET", "your-session-secret-here"),
	}

	return cfg
}

func getEnv(key string, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
