package repo

import (
	"context"
	"time"
)

// Cache缓存接口
type Cache interface {
	Put(ctx context.Context, key, value string, expire time.Duration) error
	Get(ctx context.Context, key string) (string, error)
}
