package impl

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/dgraph-io/ristretto"
)

type RistrettoCache struct {
	cache *ristretto.Cache
}

// implementation localcache

func NewRistrettoCache() (*RistrettoCache, error) {
	// ref here ANH EM: https://github.com/hypermodeinc/ristretto?tab=readme-ov-file#usage
	cache, err := ristretto.NewCache(&ristretto.Config{
		NumCounters: 1e7,     // number of keys to track frequency of (10M).
		MaxCost:     1 << 30, // maximum cost of cache (1GB).
		BufferItems: 64,      // number of keys per Get buffer.
	})
	if err != nil {
		return nil, errors.New("failed to create ristretto cache")
	}
	return &RistrettoCache{cache: cache}, nil
}

func (rc *RistrettoCache) Get(ctx context.Context, key string) (interface{}, bool) {
	return rc.cache.Get(key)
}

func (rc *RistrettoCache) Set(ctx context.Context, key string, value interface{}) bool {
	return rc.cache.Set(key, value, 1) // Cost mặc định = 1
}

func (rc *RistrettoCache) SetWithTTL(ctx context.Context, key string, value interface{}) bool {
	dataJson, _ := json.Marshal(value)
	return rc.cache.SetWithTTL(key, string(dataJson), 1, 5*time.Minute) // Cost mặc định = 1
}

func (rc *RistrettoCache) Del(ctx context.Context, key string) {
	rc.cache.Del(key)
}
