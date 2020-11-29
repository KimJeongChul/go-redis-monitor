package redis

import (
	"encoding/json"
	"strconv"
	"time"

	"github.com/KimJeongChul/go-redis-monitor/broker"
	cerror "github.com/KimJeongChul/go-redis-monitor/error"
	"github.com/KimJeongChul/go-redis-monitor/logger"
)

const packageName = "redis"

// RedisProfiler Redis server db stat
type RedisProfiler struct {
	period      int
	redisClient *RedisClient     // Redis client
	parser      *RedisInfoParser // Redis Info Parser
	broker      *broker.Broker
}

// NewRedisProfiler Create new RedisProfiler
func NewRedisProfiler(period int, redisClient *RedisClient, broker *broker.Broker) *RedisProfiler {
	rp := &RedisProfiler{
		period:      period,
		redisClient: redisClient,
		broker:      broker,
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

			logger.LogI(packageName, funcName, "# Server:", redisInfo.Server)
			logger.LogI(packageName, funcName, "# Clients:", &redisInfo.Clients)
			logger.LogI(packageName, funcName, "# Memory:", &redisInfo.Memory)
			logger.LogI(packageName, funcName, "# Persistance:", &redisInfo.Persistence)
			logger.LogI(packageName, funcName, "# Stats:", &redisInfo.Stats)
			logger.LogI(packageName, funcName, "# Replication:", &redisInfo.Replication)
			logger.LogI(packageName, funcName, "# CPU:", &redisInfo.CPU)
			logger.LogI(packageName, funcName, "# Cluster:", &redisInfo.Cluster)
			logger.LogI(packageName, funcName, "# Keyspace:", &redisInfo.Keyspace)

			connectedClient := strconv.Itoa(int(redisInfo.Clients.ConnectedClients))
			usedMmory := strconv.Itoa(int(redisInfo.Memory.UsedMemory))
			usedMemoryPeak := strconv.Itoa(int(redisInfo.Memory.UsedMemoryPeak))
			numOfDB := strconv.Itoa(len(redisInfo.Keyspace))
			totalCommandsProcessed := strconv.Itoa(int(redisInfo.Stats.TotalCommandsProcessed))
			expiredKeys := strconv.Itoa(int(redisInfo.Stats.ExpiredKeys))

			dataUpdateRedisInfo := map[string]string{
				"method":                 "updateRedisInfo",
				"connectedClient":        connectedClient,
				"usedMmory":              usedMmory,
				"usedMemoryPeak":         usedMemoryPeak,
				"totalCommandsProcessed": totalCommandsProcessed,
				"expiredKeys":            expiredKeys,
				"numOfDB":                numOfDB,
			}

			for idx, keyspace := range redisInfo.Keyspace {
				dbName := "db" + strconv.Itoa(idx)
				fieldOfKey := dbName + ":key"
				fieldOfExpires := dbName + ":expires"
				dataUpdateRedisInfo[fieldOfKey] = strconv.FormatUint(keyspace.Keys, 10)
				dataUpdateRedisInfo[fieldOfExpires] = strconv.FormatUint(keyspace.Expires, 10)
			}

			msgUpdateRedisInfo, err := json.Marshal(dataUpdateRedisInfo)
			if err == nil {
				rp.broker.Messages <- msgUpdateRedisInfo
			}
		}
	}
}

func (rp RedisProfiler) RedisLoadTest() {
	ticker := time.NewTicker(time.Duration(rp.period) * time.Second)
	for {
		select {
		case <-ticker.C:
			keyName := "TEST" + rp.getMillisTimeFormat(time.Now())
			//Set
			rp.redisClient.Set(keyName, "DONE")
			//Expire
			rp.redisClient.Expire(keyName, time.Duration(5*time.Minute))
		}
	}
}

// YYYYMMDDhhmmsslll
func (rp RedisProfiler) getMillisTimeFormat(t time.Time) string {
	timestamp := t.Format("20060102150405")
	return timestamp + strconv.Itoa(t.Nanosecond()/1000000)
}
