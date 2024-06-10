package dao

import (
	"context"
	"github.com/go-redis/redis/v8"
	"video_service/config"
	"video_service/global"
)

var Ctx = context.Background()

// InitRedis 连接redis
func InitRedis() {
	Redis := redis.NewClient(&redis.Options{
		Addr:     config.Conf.RedisConfig.Addr,
		Password: "",
		DB:       1, // 存入DB1
	})
	global.Redis = Redis
}
