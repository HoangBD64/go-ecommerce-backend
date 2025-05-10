package service

import (
	"context"
)

type IRedisCache interface {
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key string, value interface{}, expirationSeconds int) error
	Del(ctx context.Context, key string) error
	Incr(ctx context.Context, key string) (int64, error)
	Decr(ctx context.Context, key string) (int64, error)
	Exists(ctx context.Context, key string) (bool, error)

	// NEW: Distributed Lock
	WithDistributedLock(ctx context.Context, key string, ttlSeconds int, fn func(ctx context.Context) error) error
}
