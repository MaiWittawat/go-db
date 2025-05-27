package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type redisCache struct {
	redisClient *redis.Client
}


func NewRedisCache(redisClient *redis.Client) Cache {
	return &redisCache{redisClient: redisClient}
}

func (r *redisCache) Set(ctx context.Context, key string, value any, expiration time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return r.redisClient.Set(ctx, key, data, expiration).Err()
}


func (r *redisCache) Get(ctx context.Context, key string, result any) error {
	dataJson, err := r.redisClient.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil{
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


func (r *redisCache) Delete(ctx context.Context, key string) error {
	return r.redisClient.Del(ctx, key).Err()
}

