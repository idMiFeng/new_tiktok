package dao

import (
	"context"
	"github.com/go-redis/redis/v8"
	"social_service/config"
	"social_service/global"
)

var Ctx = context.Background()

// InitRedis 连接redis
func InitRedis() {
	Redis := redis.NewClient(&redis.Options{
		Addr:     config.Conf.RedisConfig.Addr,
		Password: "",
		DB:       2, // 存入DB2
	})
	global.Redis = Redis
}
