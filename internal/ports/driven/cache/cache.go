package cache

import (
	"context"
	"time"
)

// Service defines the interface for a caching service.
type Service interface {
	Get(ctx context.Context, key string) (string, error)
	GetStruct(ctx context.Context, key string, output any) error
	Set(ctx context.Context, ttl time.Duration, key string, value any) error
	Del(ctx context.Context, keys ...string) error
}
