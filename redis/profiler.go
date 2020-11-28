package redis

import (
	"fmt"
	"time"

	cerror "github.com/KimJeongChul/go-redis-monitor/error"
	"github.com/KimJeongChul/go-redis-monitor/logger"
)

const packageName = "redis"

// RedisProfiler Redis server db stat
type RedisProfiler struct {
	period      int              // 프로파일링 주기
	redisClient *RedisClient     // Redis 클라이언트
	parser      *RedisInfoParser // Redis Info Parser
}

// NewRedisProfiler Create new RedisProfiler
func NewRedisProfiler(period int, redisClient *RedisClient) *RedisProfiler {
	rp := &RedisProfiler{
		period:      period,
		redisClient: redisClient,
	}
	rp.parser = &RedisInfoParser{}
	return rp
}

func (rp RedisProfiler) Start() {
	funcName := "profiler:Start"
	ticker := time.NewTicker(time.Duration(rp.period) * time.Second)
	for {
		select {
		case <-ticker.C:
			infos, cErr := rp.redisClient.GetInfo()
			if cErr != nil {
				logger.LogE(packageName, funcName, cErr.Error())
			}

			redisInfo, cErr := rp.parser.InfoParse(infos)
			if cErr != nil {
				if cErr.Code == cerror.JSON_UNMARSHAL_ERR {
					continue
				}
			}

			fmt.Printf("# Parse result: \n")
			fmt.Printf("# Server: %+v\n", &redisInfo.Server)
			fmt.Printf("# Clients: %+v\n", &redisInfo.Clients)
			fmt.Printf("# Memory: %+v\n", &redisInfo.Memory)
			fmt.Printf("# Persistance: %+v\n", &redisInfo.Persistence)
			fmt.Printf("# Stats: %+v\n", &redisInfo.Stats)
			fmt.Printf("# Replication: %+v\n", &redisInfo.Replication)
			fmt.Printf("# CPU: %+v\n", &redisInfo.CPU)
			fmt.Printf("# Cluster: %+v\n", &redisInfo.Cluster)
			fmt.Printf("# Keyspace: %+v\n", &redisInfo.Keyspace)

			logger.LogI(packageName, funcName, "Redis info:", redisInfo.Server)

		}
	}
}
