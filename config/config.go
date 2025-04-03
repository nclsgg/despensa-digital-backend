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
		Port:          getEnv("PORT", "5310"),
	}

	return cfg
}

func getEnv(key string, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
