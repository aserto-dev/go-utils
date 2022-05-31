package cache

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/allegro/bigcache"
	"github.com/bradfitz/gomemcache/memcache"
	"github.com/coocood/freecache"
	"github.com/eko/gocache/v2/cache"
	"github.com/eko/gocache/v2/store"
	"github.com/go-redis/redis/v8"
	"github.com/rs/zerolog"
)

type CacheType int // definit aici ca enum

const (
	FreeCache CacheType = iota
	Redis
	MemCache
	BigCache
)

type CacheConfig struct {
	BigCacheConfig  bigcache.Config `json:"bigcache"`
	RedisConfig     redis.Options   `json:"redis"`
	MemCacheServers []string        `json:"memcache_servers"`
	FreeCacheSize   int             `json:"freecache_size"`
	TTL             time.Duration   `json:"ttl"`
}

func NewCache(logger *zerolog.Logger, cacheType CacheType, config CacheConfig) (cache.CacheInterface, error) {
	switch cacheType {
	case FreeCache:
		{
			freecacheStore := store.NewFreecache(freecache.NewCache(config.FreeCacheSize),
				&store.Options{Expiration: config.TTL})

			cacheManager := cache.New(freecacheStore)
			return cacheManager, nil
		}
	case MemCache:
		{
			memCacheClient := memcache.New(config.MemCacheServers...)
			memcacheStore := store.NewMemcache(
				memCacheClient,
				&store.Options{Expiration: config.TTL},
			)
			cacheManager := cache.New(memcacheStore)
			return cacheManager, nil

		}
	case Redis:
		{
			redisStore := store.NewRedis(redis.NewClient(&config.RedisConfig), &store.Options{Expiration: config.TTL})
			cacheManager := cache.New(redisStore)
			return cacheManager, nil
		}
	case BigCache:
		{
			bigcacheClient, err := bigcache.NewBigCache(config.BigCacheConfig)
			if err != nil {
				return nil, err
			}
			gocacheStore := store.NewBigcache(bigcacheClient, &store.Options{Expiration: config.TTL})
			cacheManager := cache.New(gocacheStore)
			return cacheManager, nil
		}

	default:
		{
			return nil, errors.New("cache type not supported")
		}
	}
}

func (t CacheType) String() string {
	return cacheTypeToString[t]
}

var cacheTypeToString = map[CacheType]string{
	FreeCache: "freecache",
	BigCache:  "bigcache",
	Redis:     "redis",
	MemCache:  "memcache",
}

var cacheTypeToID = map[string]CacheType{
	"bigcache":  BigCache,
	"redis":     Redis,
	"memcache":  MemCache,
	"freecache": FreeCache,
}

// MarshalJSON marshals the enum as a quoted json string
func (t CacheType) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString(`"`)
	buffer.WriteString(cacheTypeToString[t])
	buffer.WriteString(`"`)
	return buffer.Bytes(), nil
}

// UnmarshalJSON un-marshalls a quoted json string to the enum value
func (t *CacheType) UnmarshalJSON(b []byte) error {
	var j string
	err := json.Unmarshal(b, &j)
	if err != nil {
		return err
	}

	var ok bool
	*t, ok = cacheTypeToID[j]
	if !ok {
		return fmt.Errorf("'%s' is not a valid cache type", j)
	}
	return nil
}
