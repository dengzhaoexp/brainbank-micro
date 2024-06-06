package cache

import (
	"fmt"
	"github.com/go-redis/redis"
	"user/config"
	userLogger "user/pkg/utils/logger"
)

var _redisClient *redis.Client

func InitRedis() {
	rConfig := config.Config.Redis
	userLogger.LogrusObj.Error("Loading Redis Configuration Successfully.")

	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", rConfig.Host, rConfig.Port),
		Password: rConfig.Password,
		DB:       rConfig.RedisDbName,
	})
	if _, err := client.Ping().Result(); err != nil {
		userLogger.LogrusObj.Error("Error while pinging Redis server:", err)
		panic(err)
	}

	userLogger.LogrusObj.Error("Client for redis created successfully.")
	_redisClient = client
}

func GetRedisClient() *redis.Client {
	return _redisClient
}
