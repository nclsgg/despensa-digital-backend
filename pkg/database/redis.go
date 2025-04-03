package database

import (
	"context"
	"log"
	"time"

	"github.com/nclsgg/despensa-digital/backend/config"
	"github.com/redis/go-redis/v9"
)

func ConnectRedis(cfg *config.Config) *redis.Client {
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
			log.Printf("Error connecting to Redis: %v\n", err)
			log.Printf("Retrying in %v...\n", delay)
			time.Sleep(delay)
			continue
		}

		log.Println("Connected to Redis")
		return client
	}

	log.Fatalf("Failed to connect to Redis after %d attempts", maxAttempts)
	return nil
}
