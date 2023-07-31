package dao

import (
	"context"
	"dragonsss.cn/evn_project/config"
	"github.com/go-redis/redis/v8"
	"time"
)

var Rc *RedisCache

// RedisCache redis缓存
type RedisCache struct {
	rdb *redis.Client
}

func init() {
	//连接redis 客户端
	rdb := redis.NewClient(config.C.ReadRedisConfig()) //读取配置文件
	Rc = &RedisCache{
		rdb: rdb,
	}
}

func (rc *RedisCache) Put(ctx context.Context, key, value string, expire time.Duration) error {
	err := rc.rdb.Set(ctx, key, value, expire).Err()
	return err
}

func (rc *RedisCache) Get(ctx context.Context, key string) (string, error) {
	result, err := rc.rdb.Get(ctx, key).Result()
	return result, err
}
