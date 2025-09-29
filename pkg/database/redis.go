package database

import (
	"context"
	"log"
	"time"

	"github.com/nclsgg/despensa-digital/backend/config"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

func ConnectRedis(cfg *config.Config) (result0 *redis.Client) {
	__logParams := map[string]any{"cfg": cfg}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "ConnectRedis"), zap.Any("result", result0), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "ConnectRedis"), zap.Any("params", __logParams))
	log.Println("Connecting to Redis...")
	ctx := context.Background()

	maxAttempts := 3
	delay := 2 * time.Second

	for attempt := 1; attempt <= maxAttempts; attempt++ {
		client := redis.NewClient(&redis.Options{
			Addr:     cfg.RedisURL,
			Username: cfg.RedisUsername,
			Password: cfg.RedisPassword,
			DB:       0,
		})

		_, err := client.Ping(ctx).Result()
		if err != nil {
			zap.L().Error("function.error", zap.String("func", "ConnectRedis"), zap.Error(err), zap.Any("params", __logParams))
			log.Printf("Error connecting to Redis: %v\n", err)
			log.Printf("Retrying in %v...\n", delay)
			time.Sleep(delay)
			continue
		}

		log.Println("Connected to Redis")
		result0 = client
		return
	}

	log.Fatalf("Failed to connect to Redis after %d attempts", maxAttempts)
	result0 = nil
	return
}
