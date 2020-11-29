# go-redis-monitor
Redis monitor displays information about the Redis db stats in real time.

![image](https://user-images.githubusercontent.com/10591350/100537620-44154200-326d-11eb-91e6-a6aea8963fe8.png)

#### Redis info
```bash
$ redis-cli -p 7001 -a passw0rd info
```
![image](https://user-images.githubusercontent.com/10591350/100492751-83ab3380-3172-11eb-84c2-28ec4b111055.png)

#### Parse result
![image](https://user-images.githubusercontent.com/10591350/100492810-61fe7c00-3173-11eb-8410-895c6db76a67.png)

### Start
```shell
go build
./go-redis-monitor
```

### Config file
config.json
```json
{
    "port": {HTTP_SERVER_PORT},
    "period": {MONITORING_PERIOD},
    "redisAddr": {REDIS_SERVER_ADDR},
    "redisPort": {REDIS_SERVER_PORT},
    "redisDB": {REDIS_DB_NAME},
    "redisPassword": {REDIS_PASSWORD},
    "logPath": {LOG_FILE_PATH}
}
```