package cache

import (
	"errors"
	"time"

	"github.com/dgraph-io/ristretto"
)

// ErrValNotSet error returned when value not set
var ErrValNotSet = errors.New("unable to set value in cache")

// Cache in-memory cache
type Cache struct {
	cache *ristretto.Cache
}

// New creates a new cache
func New(cache *ristretto.Cache) *Cache {
	return &Cache{
		cache: cache,
	}
}

// Set sets the cache value
func (c *Cache) Set(key string, val interface{}, ttl time.Duration) error {
	if !c.cache.SetWithTTL(key, val, 0, ttl) {
		return ErrValNotSet
	}

	c.cache.Wait()

	return nil
}

// Get fetches cached value
func (c *Cache) Get(key string) (interface{}, error) {
	v, ok := c.cache.Get(key)
	if !ok {
		return nil, nil
	}

	return v, nil
}

// Clear clears cached value
func (c *Cache) Clear() {
	c.cache.Clear()
}
