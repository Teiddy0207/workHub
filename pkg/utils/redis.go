package utils

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

var (
	RedisClient *redis.Client
	Ctx         = context.Background()
)

func InitRedis(cfgHost string, cfgPort int, cfgPassword string, cfgDB int, poolSize int, readTimeout int, writeTimeout int, dialTimeout int, timeout int) error {
	RedisClient = redis.NewClient(&redis.Options{
		Addr:         fmt.Sprintf("%s:%d", cfgHost, cfgPort),
		Password:     cfgPassword,
		DB:           cfgDB,
		PoolSize:     poolSize,
		ReadTimeout:  time.Duration(readTimeout) * time.Second,
		WriteTimeout: time.Duration(writeTimeout) * time.Second,
		DialTimeout:  time.Duration(dialTimeout) * time.Second,
		ContextTimeoutEnabled: true,
	})

	if err := RedisClient.Ping(Ctx).Err(); err != nil {
		return fmt.Errorf("failed to connect to Redis: %w", err)
	}
	
	fmt.Println("Connected to Redis successfully")
	return nil
}

func SaveToRedis(ctx context.Context, client *redis.Client, key string, value interface{}, expiration time.Duration) error {
	err := client.Set(ctx, key, value, expiration).Err()
	if err != nil {
		return fmt.Errorf("failed to save key '%s' to redis: %w", key, err)
	}
	return nil
}

func GetFromRedis(ctx context.Context, client *redis.Client, key string) (string, error) {
	val, err := client.Get(ctx, key).Result()
	if err != nil {
		return "", fmt.Errorf("failed to get key '%s' from redis: %w", key, err)
	}
	return val, nil
}
