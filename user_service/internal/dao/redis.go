package dao

import (
	"context"
	"github.com/go-redis/redis/v8"
	"user_service/config"
	"user_service/global"
)

var Ctx = context.Background()

// InitRedis 连接redis
func InitRedis() {
	Redis := redis.NewClient(&redis.Options{
		Addr:     config.Conf.RedisConfig.Addr,
		Password: "",
		DB:       0, // 存入DB0
	})
	global.Redis = Redis
}
