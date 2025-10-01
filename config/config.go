package config

import (
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

type Config struct {
	DatabaseURL   string
	JWTSecret     string
	JWTExpiration string
	JWTIssuer     string
	JWTAudience   string
	Port          string
	CorsOrigin    string

	// OAuth Config
	GoogleClientID     string
	GoogleClientSecret string
	GoogleCallbackURL  string
	FrontendURL        string
	SessionSecret      string
}

func LoadConfig() (result0 *Config) {
	__logParams := map[string]any{}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "LoadConfig"), zap.Any("result", result0), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "LoadConfig"), zap.Any("params", __logParams))
	_ = godotenv.Load()

	cfg := &Config{
		DatabaseURL:   os.Getenv("DATABASE_URL"),
		JWTSecret:     os.Getenv("JWT_SECRET"),
		JWTExpiration: os.Getenv("JWT_EXPIRATION"),
		JWTIssuer:     os.Getenv("JWT_ISSUER"),
		JWTAudience:   os.Getenv("JWT_AUDIENCE"),
		Port:          getEnv("PORT", "3030"),
		CorsOrigin:    getEnv("CORS_ORIGIN", "http://localhost:3000"),

		// OAuth Config
		GoogleClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		GoogleClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		GoogleCallbackURL:  getEnv("GOOGLE_CALLBACK_URL", "http://localhost:3030/auth/oauth/google/callback"),
		FrontendURL:        getEnv("FRONTEND_URL", "http://localhost:3000"),
		SessionSecret:      getEnv("SESSION_SECRET", "your-session-secret-here"),
	}
	result0 = cfg
	return
}

func getEnv(key string, fallback string) (result0 string) {
	__logParams := map[string]any{"key": key, "fallback": fallback}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "getEnv"), zap.Any("result", result0), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "getEnv"), zap.Any("params", __logParams))
	if value, exists := os.LookupEnv(key); exists {
		result0 = value
		return
	}
	result0 = fallback
	return
}
