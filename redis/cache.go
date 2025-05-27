package redis

import (
	"context"
	"time"
)

type Cache interface {
	Set(ctx context.Context, key string, value any, expiration time.Duration) error
	Get(ctx context.Context, key string, result any) error
	Delete(ctx context.Context, key string) error
}