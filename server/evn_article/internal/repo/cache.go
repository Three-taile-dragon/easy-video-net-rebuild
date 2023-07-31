package repo

import (
	"context"
	"time"
)

// Cache 缓存接口	方便后续更换缓存存储软件 redis mysql mongo...
type Cache interface {
	Put(ctx context.Context, key, value string, expire time.Duration) error
	Get(ctx context.Context, key string) (string, error)
}
