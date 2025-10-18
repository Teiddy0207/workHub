package utils

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

func SaveToRedis(ctx context.Context, client *redis.Client, key string, value interface{}, expiration time.Duration) error {
	err := client.Set(ctx, key, value, expiration).Err()
	if err != nil {
		return fmt.Errorf("failed to save key '%s' to redis: %w", key, err)
	}
	return nil
}