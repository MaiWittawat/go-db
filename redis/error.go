package redis

import "errors"

var (
	ErrCacheMiss = errors.New("cache miss")
)