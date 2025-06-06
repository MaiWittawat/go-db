package cache

import (
	"context"
	"errors"
	"fmt"
	"time"
)

var (
	ErrCacheMiss = errors.New("cache miss")
)

type Cache interface {
	Set(ctx context.Context, key string, value any, expiration time.Duration) error
	Get(ctx context.Context, key string, result any) error
	Delete(ctx context.Context, key string) error
}

type KeyGenerator struct {
	prefix string	
}

func NewKeyGenerator(prefix string) *KeyGenerator {
	return &KeyGenerator{prefix: prefix}
}

func (k *KeyGenerator) KeyList() string {
	return k.prefix
}

func (k *KeyGenerator) KeyID(id string) string {
	return fmt.Sprintf("%s:%s", k.prefix, id)
}

func (k *KeyGenerator) KeyField(field string, value string) string {
	return fmt.Sprintf("%s:%s:%s", k.prefix, field, value)
}