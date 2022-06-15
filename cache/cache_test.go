package cache

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/alecthomas/assert"
	"github.com/allegro/bigcache"
	"github.com/rs/zerolog"
)

func TestBigCache(t *testing.T) {
	mycache, err := NewCache(&zerolog.Logger{},
		BigCache,
		CacheConfig{
			BigCacheConfig: bigcache.DefaultConfig(time.Duration(5 * time.Minute)),
			TTL:            time.Duration(5 * time.Minute),
		})
	assert.NoError(t, err)
	err = mycache.Set(context.Background(), "test", "testValue", nil)
	assert.NoError(t, err)
	value, err := mycache.Get(context.Background(), "test")
	assert.NoError(t, err)
	assert.Equal(t, fmt.Sprintf("%s", value), "testValue")
}

func TestFreeCache(t *testing.T) {
	mycache, err := NewCache(&zerolog.Logger{},
		FreeCache,
		CacheConfig{
			FreeCacheSize: 1000,
			TTL:           time.Duration(5 * time.Minute),
		})
	assert.NoError(t, err)
	// freeCache allows setting only byte array values
	err = mycache.Set(context.Background(), "test", []byte("testValue"), nil)
	assert.NoError(t, err)
	value, err := mycache.Get(context.Background(), "test")
	assert.NoError(t, err)
	assert.Equal(t, fmt.Sprintf("%s", value), "testValue")
}
