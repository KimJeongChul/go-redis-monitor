package redis

import (
	"context"
	"time"

	cerror "github.com/KimJeongChul/go-redis-monitor/error"
	"github.com/KimJeongChul/go-redis-monitor/logger"
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

// Set SET command
func (rc *RedisClient) Set(key string, value string) (string, *cerror.CError) {
	result, err := rc.rdb.Set(rc.ctx, key, value, 0).Result()
	if err != nil {
		cErr := cerror.NewCError(cerror.REDIS_DB_ERR, "redis Set error:"+err.Error())
		return "", cErr
	} else {
		logger.LogI(packageName, "Set", "key:"+key+" value:"+value+" done!")
	}
	return result, nil
}

// Expire EXPIRE command
func (rc *RedisClient) Expire(key string, period time.Duration) (bool, *cerror.CError) {
	result, err := rc.rdb.Expire(rc.ctx, key, period).Result()
	if err != nil {
		cErr := cerror.NewCError(cerror.REDIS_DB_ERR, "redis Expire error:"+err.Error())
		return result, cErr
	} else {
		logger.LogI(packageName, "Expire", "key:"+key+" period:"+period.String()+" done!")
	}
	return result, nil
}
