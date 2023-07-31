package dao

import (
	"context"
	"dragonsss.cn/evn_user/config"
	"github.com/go-redis/redis/v8"
	"time"
)

var Rc *RedisCache

type RedisCache struct {
	rdb *redis.Client
}

func init() {
	//连接redis客户端
	rdb := redis.NewClient(config.C.ReadRedisConfig())
	Rc = &RedisCache{
		rdb: rdb,
	}
}

// Put 实现缓存存入接口
func (rc *RedisCache) Put(ctx context.Context, key, value string, expire time.Duration) error {
	err := rc.rdb.Set(ctx, key, value, expire).Err()
	return err
}

// Get 实现缓存查询接口
func (rc *RedisCache) Get(ctx context.Context, key string) (string, error) {
	result, err := rc.rdb.Get(ctx, key).Result()
	return result, err
}
