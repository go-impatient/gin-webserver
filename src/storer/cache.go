package storer

import (
	"time"

	"github.com/go-redis/redis"
	cache "github.com/patrickmn/go-cache"
)

//CacheStore represents the cache store
type CacheStore struct {
	DefaultTimeout time.Duration
	Type           uint8
	redis          *redis.Client
	gocache        *cache.Cache
}

// NewCacheStore initializes the cache store
func NewCacheStore() *CacheStore {
	return &CacheStore{

	}
}

//Set
func (c *CacheStore) Set(cKey string) error {

	return nil
}

//Get
func (c *CacheStore) Get(cKey string) error {
	return nil
}
