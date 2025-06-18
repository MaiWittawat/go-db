package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

type cacheService struct {
	redisClient *redis.Client
}

func InitRedisClient(addr string, password string) *redis.Client {
	redisAddr := addr
	redisPass := password
	rdb := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: redisPass,
	})

	// check connect
	_, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		log.Panic("Failed to connect to Redis:", err)
	}
	return rdb
}

// ------------------------ Constructor ------------------------
func NewCacheService(redisClient *redis.Client) Cache {
	return &cacheService{redisClient: redisClient}
}

// ------------------------ Method Basic Set, Get, Del ------------------------
func (s *cacheService) Set(ctx context.Context, key string, value any, expiration time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return s.redisClient.Set(ctx, key, data, expiration).Err()
}

func (s *cacheService) Get(ctx context.Context, key string, result any) error {
	dataJson, err := s.redisClient.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return ErrCacheMiss
		}
		return fmt.Errorf("fail to get from cache: %w", err)
	}

	// Unmarshall จะเปลี่ยนค่ากลับเป็น ตัวเเปร เป้าหมายรับเฉพาะค่า pointer เท่านั้น
	if err := json.Unmarshal([]byte(dataJson), result); err != nil {
		return fmt.Errorf("fail to unmarshal cache data: %w", err)
	}
	return nil
}

func (s *cacheService) Delete(ctx context.Context, key string) error {
	return s.redisClient.Del(ctx, key).Err()
}
