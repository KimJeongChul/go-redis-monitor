package redis

import (
	"encoding/json"
	"strconv"
	"strings"
	"unicode/utf8"

	cerror "github.com/KimJeongChul/go-redis-monitor/error"
	"github.com/KimJeongChul/go-redis-monitor/logger"
)

// RedisInfoParser
type RedisInfoParser struct {
}

// InfoParse Parse Redis INFO stat
func (rip RedisInfoParser) InfoParse(info string) (*RedisInfo, *cerror.CError) {
	redisInfo, cErr := rip.Parse(info)
	if cErr != nil {
		return nil, cErr
	}
	return &redisInfo, nil
}

type RedisInfo struct {
	Server      Server      `json:"server"`
	Clients     Client      `json:"clients"`
	Memory      Memory      `json:"memory"`
	Persistence Persistence `json:"persistence"`
	Stats       Stats       `json:"stats"`
	Replication Replication `json:"replication"`
	CPU         CPU         `json:"cpu"`
	Cluster     Cluster     `json:"cluster"`
	Keyspace    []Keyspace  `json:"keyspace"`
}

type Server struct {
	RedisVersion    string `json:"redis_version"`
	RedisGitSha1    string `json:"redis_git_sha1"`
	RedisGitDirty   uint8  `json:"redis_git_dirty"`
	RedisBuildID    string `json:"redis_build_id"`
	RedisMode       string `json:"redis_mode"`
	OS              string `json:"os"`
	ArchBits        uint8  `json:"arch_bits"`
	MultiplexingAPI string `json:"multiplexing_api"`
	GCCVersion      string `json:"gcc_version"`
	ProcessID       uint16 `json:"process_id"`
	RunID           string `json:"run_id"`
	TCPPort         uint16 `json:"tcp_port"`
	UptimeInSeconds uint64 `json:"uptime_in_seconds"`
	UptimeInDays    uint64 `json:"uptime_in_days"`
	HZ              uint64 `json:"hz"`
	LRUClock        int64  `json:"lru_clock"`
	Executable      string `json:"executable"`
	ConfigFile      string `json:"config_file"`
}

type Client struct {
	ConnectedClients        int64 `json:"connected_clients"`
	ClientLongestOutputList int64 `json:"client_longest_output_list"`
	ClientBiggestInputBuf   int64 `json:"client_biggest_input_buf"`
	BlockedClients          int64 `json:"blocked_clients"`
}

type Memory struct {
	UsedMemory             int64   `json:"used_memory"`
	UsedMemoryHuman        string  `json:"used_memory_human"`
	UsedMemoryRss          int64   `json:"used_memory_rss"`
	UsedMemoryRssHuman     string  `json:"used_memory_rss_human"`
	UsedMemoryPeak         int64   `json:"used_memory_peak"`
	UsedMemoryPeakHuman    string  `json:"used_memory_peak_human"`
	TotalSystemMemory      int64   `json:"total_system_memory"`
	TotalSystemMemoryHuman string  `json:"total_system_memory_human"`
	UsedMemoryLua          int64   `json:"used_memory_lua"`
	UsedMemoryLuaHuman     string  `json:"used_memory_lua_human"`
	Maxmemory              int64   `json:"maxmemory"`
	MaxmemoryHuman         string  `json:"maxmemory_human"`
	MaxmemoryPolicy        string  `json:"maxmemory_policy"`
	MemFragmentationRatio  float64 `json:"mem_fragmentation_ratio"`
	MemAllocator           string  `json:"mem_allocator"`
}

type Persistence struct {
	Loading                  int64  `json:"loading"`
	RdbChangesSinceLastSave  int64  `json:"rdb_changes_since_last_save"`
	RdbBgsaveInProgress      int64  `json:"rdb_bgsave_in_progress"`
	RdbLastSaveTime          int64  `json:"rdb_last_save_time"`
	RdbLastBgsaveStatus      string `json:"rdb_last_bgsave_status"`
	RdbLastBgsaveTimeSec     int64  `json:"rdb_last_bgsave_time_sec"`
	RdbCurrentBgsaveTimeSec  int64  `json:"rdb_current_bgsave_time_sec"`
	AofEnabled               int64  `json:"aof_enabled"`
	AofRewriteInProgress     int64  `json:"aof_rewrite_in_progress"`
	AofRewriteScheduled      int64  `json:"aof_rewrite_scheduled"`
	AofLastRewriteTimeSec    int64  `json:"aof_last_rewrite_time_sec"`
	AofCurrentRewriteTimeSec int64  `json:"aof_current_rewrite_time_sec"`
	AofLastBgrewriteStatus   string `json:"aof_last_bgrewrite_status"`
	AofLastWriteStatus       string `json:"aof_last_write_status"`
}

type Stats struct {
	TotalConnectionsReceived int64   `json:"total_connections_received"`
	TotalCommandsProcessed   int64   `json:"total_commands_processed"`
	InstantaneousOpsPerSec   int64   `json:"instantaneous_ops_per_sec"`
	TotalNetInputBytes       int64   `json:"total_net_input_bytes"`
	TotalNetOutputBytes      int64   `json:"total_net_output_bytes"`
	InstantaneousInputKbps   float64 `json:"instantaneous_input_kbps"`
	InstantaneousOutputKbps  float64 `json:"instantaneous_output_kbps"`
	RejectedConnections      int64   `json:"rejected_connections"`
	SyncFull                 int64   `json:"sync_full"`
	SyncPartialOk            int64   `json:"sync_partial_ok"`
	SyncPartialErr           int64   `json:"sync_partial_err"`
	ExpiredKeys              int64   `json:"expired_keys"`
	EvictedKeys              int64   `json:"evicted_keys"`
	KeyspaceHits             int64   `json:"keyspace_hits"`
	KeyspaceMisses           int64   `json:"keyspace_misses"`
	PubsubChannels           int64   `json:"pubsub_channels"`
	PubsubPatterns           int64   `json:"pubsub_patterns"`
	LatestForkUsec           int64   `json:"latest_fork_usec"`
	MigrateCachedSockets     int64   `json:"migrate_cached_sockets"`
}

type Cluster struct {
	ClusterEnabled int64 `json:"cluster_enabled"`
}

type Replication struct {
	Role                       string             `json:"role"`
	ConnectedSlaves            uint64             `json:"connected_slaves"`
	Slaves                     []ReplicationSlave `json:"slave"`
	MasterReplOffset           int64              `json:"master_repl_offset"`
	ReplBacklogActive          int64              `json:"repl_backlog_active"`
	ReplBacklogSize            uint64             `json:"repl_backlog_size"`
	ReplBacklogFirstByteOffset uint64             `json:"repl_backlog_first_byte_offset"`
	ReplBacklogHistLen         uint64             `json:"repl_backlog_histlen"`
}

type ReplicationSlave struct {
	ID     int64  `json:"id,omitempty" index:"true"`
	IP     string `json:"ip"`
	Port   uint16 `json:"port"`
	State  string `json:"state"`
	Offset int64  `json:"offset"`
	Lag    int64  `json:"lag"`
}

type CPU struct {
	UsedCPUSys          float64 `json:"used_cpu_sys"`
	UsedCPUUser         float64 `json:"used_cpu_user"`
	UsedCPUSysChildren  float64 `json:"used_cpu_sys_children"`
	UsedCPUUserChildren float64 `json:"used_cpu_user_children"`
}

type Keyspace struct {
	DB      uint64 `json:"db,omitempty" index:"true"`
	Keys    uint64 `json:"keys"`
	Expires uint64 `json:"expires"`
	AvgTTL  uint64 `json:"avg_ttl"`
}

func (rip RedisInfoParser) Parse(content string) (ri RedisInfo, cErr *cerror.CError) {
	funcName := "parser:Parse"

	ri = RedisInfo{}

	// Carriage return 제거
	content = strings.Replace(content, "\r", "", -1)
	lines := strings.Split(content, "\n")

	redisInfoValues := make(map[string][]string, 0)
	sectionName := ""

	//Preprocessing
	for _, line := range lines {
		if line == "" {
			continue
		} else if strings.HasPrefix(line, "# ") {
			sectionName = strings.ToLower(line[2:])
			redisInfoValues[sectionName] = make([]string, 0)
		} else {
			redisInfoValues[sectionName] = append(redisInfoValues[sectionName], line)
		}
	}

	var keyspaces []Keyspace
	for sectionName, _ := range redisInfoValues {
		jsonString := "{"
		for _, v := range redisInfoValues[sectionName] {

			splV := strings.Split(v, ":")
			key := "\"" + splV[0] + "\""
			value := splV[1]

			if sectionName == "keyspace" {
				var keyspace Keyspace
				indexDB := splV[0][2:]
				uIndexDB, _ := strconv.ParseUint(indexDB, 10, 64)
				keyspace.DB = uIndexDB
				dbInfos := strings.Split(splV[1], ",")
				for _, dbInfo := range dbInfos {
					splDBInfo := strings.Split(dbInfo, "=")
					uDBInfoValue, _ := strconv.ParseUint(splDBInfo[1], 10, 64)
					switch splDBInfo[0] {
					case "keys":
						keyspace.Keys = uDBInfoValue
					case "expires":
						keyspace.Expires = uDBInfoValue
					case "avg_ttl":
						keyspace.AvgTTL = uDBInfoValue
					}
				}
				keyspaces = append(keyspaces, keyspace)
			}
			if intV, err := strconv.Atoi(value); err != nil {
				if _, err := strconv.ParseFloat(value, 64); err != nil {
					value = "\"" + value + "\""
				}
			} else {
				if intV == 0 && len(value) > 1 {
					value = "\"" + value + "\""
				}
			}
			jsonString += key + ":" + value + ","
		}
		jsonString = trimLastChar(jsonString)
		jsonString += "}"

		switch sectionName {
		case "server":
			var server Server
			err := json.Unmarshal([]byte(jsonString), &server)
			if err != nil {
				logger.LogE(packageName, funcName, "section name:server err:"+err.Error())
				cErr = cerror.NewCError(cerror.JSON_UNMARSHAL_ERR, "section name:server err:"+err.Error())
			} else {
				ri.Server = server
			}
		case "clients":
			var client Client
			err := json.Unmarshal([]byte(jsonString), &client)
			if err != nil {
				logger.LogE(packageName, funcName, "section name:clients err:"+err.Error())
				cErr = cerror.NewCError(cerror.JSON_UNMARSHAL_ERR, "section name:clients err:"+err.Error())
			} else {
				ri.Clients = client
			}
		case "memory":
			var memory Memory
			err := json.Unmarshal([]byte(jsonString), &memory)
			if err != nil {
				logger.LogE(packageName, funcName, "section name:memory err:"+err.Error())
				cErr = cerror.NewCError(cerror.JSON_UNMARSHAL_ERR, "section name:memory err:"+err.Error())
			} else {
				ri.Memory = memory
			}
		case "persistence":
			var persistence Persistence
			err := json.Unmarshal([]byte(jsonString), &persistence)
			if err != nil {
				logger.LogE(packageName, funcName, "section name:persistance err:"+err.Error())
				cErr = cerror.NewCError(cerror.JSON_UNMARSHAL_ERR, "section name:persistance err:"+err.Error())
			} else {
				ri.Persistence = persistence
			}
		case "stats":
			var stats Stats
			err := json.Unmarshal([]byte(jsonString), &stats)
			if err != nil {
				logger.LogE(packageName, funcName, "section name:stats err:"+err.Error())
				cErr = cerror.NewCError(cerror.JSON_UNMARSHAL_ERR, "section name:stats err:"+err.Error())
			} else {
				ri.Stats = stats
			}
		case "cpu":
			var cpu CPU
			err := json.Unmarshal([]byte(jsonString), &cpu)
			if err != nil {
				logger.LogE(packageName, funcName, "section name:cpu err:"+err.Error())
				cErr = cerror.NewCError(cerror.JSON_UNMARSHAL_ERR, "section name:cpu err:"+err.Error())
			} else {
				ri.CPU = cpu
			}
		case "replication":
			var replication Replication
			err := json.Unmarshal([]byte(jsonString), &replication)
			if err != nil {
				logger.LogE(packageName, funcName, "section name:replication err:"+err.Error())
				cErr = cerror.NewCError(cerror.JSON_UNMARSHAL_ERR, "section name:replication err:"+err.Error())
			} else {
				ri.Replication = replication
			}
		case "cluster":
			var cluster Cluster
			err := json.Unmarshal([]byte(jsonString), &cluster)
			if err != nil {
				logger.LogE(packageName, funcName, "section name:cluster err:"+err.Error())
				cErr = cerror.NewCError(cerror.JSON_UNMARSHAL_ERR, "section name:cluster err:"+err.Error())
			} else {
				ri.Cluster = cluster
			}
		}
	}

	ri.Keyspace = keyspaces
	return
}

// trimLastChar Delete last character from string
func trimLastChar(s string) string {
	r, size := utf8.DecodeLastRuneInString(s)
	if r == utf8.RuneError && (size == 0 || size == 1) {
		size = 0
	}
	return s[:len(s)-size]
}
