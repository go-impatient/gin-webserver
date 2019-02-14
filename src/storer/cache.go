package storer

import (
	"fmt"
	"time"

	"github.com/go-redis/redis"
	"github.com/moocss/go-webserver/src/config"
	"github.com/patrickmn/go-cache"
	"bytes"
	"encoding/json"
)

//CacheStore represents the cache store
type CacheStore struct {
	DefaultTimeout time.Duration
	Type           uint8
	redis          *redis.Client
	gocache        *cache.Cache
}

const (
	typeLocal = 0
	typeRedis = 1
	typeNone  = 2
)

//InitCacheStore initializes the cache store
func InitCacheStore(cfg *config.ConfigCache) *CacheStore {
	addr := fmt.Sprintf("%s:%s", cfg.Redis.Host, cfg.Redis.Port)
	cs := &CacheStore{
		DefaultTimeout: time.Duration(cfg.Timeout) * time.Second,
	}

	switch cfg.Type {
	case "local":
		cs.Type = typeLocal
		cs.gocache = cache.New(cs.DefaultTimeout, 4*time.Minute)
	case "redis":
		cs.Type = typeRedis
		cs.redis = redis.NewClient(&redis.Options{
			Addr:     addr,
			Password: cfg.Redis.Password,
			DB:       cfg.Redis.DB,
		})
	case "none":
		cs.Type = typeNone
	default:
		cs.Type = typeLocal
	}

	return cs
}

// Set stores
func (c *CacheStore) Set(key string, data string) error {
	switch c.Type {
		case typeLocal:
			c.gocache.Set(key, data, cache.DefaultExpiration)
		case typeRedis:
			op := &bytes.Buffer{}
			enc := json.NewEncoder(op)
			enc.SetEscapeHTML(false)
			err := enc.Encode(data)
			if err != nil {
				return err
			}
			err = c.redis.SetNX(key, string(op.Bytes()), c.DefaultTimeout).Err()
			if err != nil {
				return err
			}
		case typeNone:
		default:
	}
	return nil
}

// Get gets
func (c *CacheStore) Get(key string) (string, error) {
	switch c.Type {
		case typeLocal:
			if x, found := c.gocache.Get(key); found {
				foo := x.(string)
				return foo, nil
			}
			return "", nil
		case typeRedis:
			val, err := c.redis.Get(key).Result()
			if err == redis.Nil {
				return "", nil
			} else if err != nil {
				return "", err
			} else {
				return val, nil
			}
		case typeNone:
			return "", nil
		default:
			return "", nil
	}
}

