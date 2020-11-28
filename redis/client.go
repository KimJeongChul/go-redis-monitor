package redis

import (
	"context"

	cerror "github.com/KimJeongChul/go-redis-monitor/error"
	"github.com/go-redis/redis/v8"
)

// RedisSentinelClient RedisClient
type RedisClient struct {
	rdb *redis.Client
	ctx context.Context
}

// NewRedisSentinelClient Create new reids client
func NewRedisClient(ctx context.Context, addr string, port string, password string, db string) *RedisClient {
	redisServerAddr := addr + ":" + port
	rdb := redis.NewClient(&redis.Options{
		Addr:     redisServerAddr,
		Password: password,
		DB:       0,
	})

	rdb.Ping(ctx)

	rc := &RedisClient{
		rdb: rdb,
		ctx: ctx,
	}

	return rc
}

// GetInfo Get Info
func (rc *RedisClient) GetInfo() (string, *cerror.CError) {
	result, err := rc.rdb.Info(rc.ctx).Result()
	if err != nil {
		cErr := cerror.NewCError(cerror.REDIS_DB_ERR, "redis info error:"+err.Error())
		return "", cErr
	}
	return result, nil
}
