package cache

import (
	"context"
	"fmt"

	errors "github.com/webitel/engine/model"
	"github.com/webitel/engine/utils"
)

type MemoryCache struct {
	lruCache utils.ObjectCache
}

type MemoryCacheConfig struct {
	// Size of the cache
	Size int
	// Default expire in seconds
	DefaultExpiry int64
}

func NewMemoryCache(conf *MemoryCacheConfig) *MemoryCache {
	return &MemoryCache{lruCache: utils.NewLruWithParams(conf.Size, "memoryCache", conf.DefaultExpiry, "")}
}

func (m *MemoryCache) Get(ctx context.Context, key string) (*CacheValue, errors.AppError) {
	value, ok := m.lruCache.Get(key)
	if !ok {
		return nil, errors.NewInternalError("cache.memory_cache.get", fmt.Sprintf("unable to find value by key %s", key))
	}
	return NewCacheValue(value)
}

func (m *MemoryCache) Set(ctx context.Context, key string, value any, expiresAfterSecs int64) errors.AppError {
	m.lruCache.AddWithExpiresInSecs(key, value, expiresAfterSecs)
	return nil
}
func (m *MemoryCache) Delete(ctx context.Context, key string) errors.AppError {
	m.lruCache.Remove(key)
	return nil
}

func (m *MemoryCache) IsValid() errors.AppError {
	if m.lruCache == nil {
		return errors.NewInternalError("cache.memory_cache", "lru cache client is not declared")
	}
	return nil
}
