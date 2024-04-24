package cache

import (
	"context"
	"encoding/json"
	"time"

	"github.com/go-redis/redis/v8"

	"reprocess-gui/internal/errors"
)

type redisCache struct {
	reader *redis.Client
	writer *redis.Client
}

// NewRedisService creates a new instance of the Redis cache service.
func NewRedisService(reader, writer *redis.Client) Service {
	return &redisCache{
		reader: reader,
		writer: writer,
	}
}

// Get retrieves the value associated with the provided key from the Redis cache.
func (c *redisCache) Get(ctx context.Context, key string) (string, error) {
	if c.reader == nil {
		return "", errors.ErrRedisReaderClientNotInitialized
	}

	val, err := c.reader.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", errors.ErrCacheNotFound
	}
	if err != nil {
		return "", err
	}

	return val, nil
}

// GetStruct retrieves the value associated with the provided key from the Redis cache and unmarshals it into the provided output.
func (c *redisCache) GetStruct(ctx context.Context, key string, output any) error {
	res, err := c.Get(ctx, key)
	if err != nil {
		return err
	}

	return json.Unmarshal([]byte(res), output)
}

// Set sets the value associated with the provided key in the Redis cache.
func (c *redisCache) Set(ctx context.Context, ttl time.Duration, key string, value any) error {
	if c.writer == nil {
		return errors.ErrRedisWriterClientNotInitialized
	}

	return c.writer.Set(ctx, key, value, ttl).Err()
}

// Del deletes the values associated with the provided keys from the Redis cache.
func (c *redisCache) Del(ctx context.Context, keys ...string) error {
	if c.writer == nil {
		return errors.ErrRedisWriterClientNotInitialized
	}

	return c.writer.Del(ctx, keys...).Err()
}
