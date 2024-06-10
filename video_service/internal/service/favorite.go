package service

import (
	"fmt"
	"strconv"
	"user_service/logger"
	"video_service/global"
	"video_service/internal/dao"
	"video_service/internal/model"
)

// IsFavorite 检查用户是否点赞某个视频
func IsFavorite(userId, videoId uint) bool {
	key := fmt.Sprintf("video:user_fav:%s", strconv.FormatInt(int64(userId), 10))
	videoIdStr := strconv.FormatUint(uint64(videoId), 10)

	// 先检查 Redis 中是否存在该 key
	exists, err := global.Redis.Exists(dao.Ctx, key).Result()
	if err != nil {
		logger.Log.Error("缓存错误：", err)
	}

	if exists > 0 {
		// 如果存在，检查视频 ID 是否在集合中
		isMember, err := global.Redis.SIsMember(dao.Ctx, key, videoIdStr).Result()
		if err != nil {
			logger.Log.Error("缓存错误：", err)
			return false

		}
		return isMember
	} else {
		// 如果不存在，从 MySQL 查询
		isFavorite, err := model.GetFavoriteInstance().IsFavoriteVideo(userId, videoId)
		if err != nil {
			logger.Log.Error("缓存错误：", err)
			return false
		}

		// 将结果写回 Redis
		if isFavorite {
			_, err := global.Redis.SAdd(dao.Ctx, key, videoIdStr).Result()
			if err != nil {
				logger.Log.Error("缓存错误：", err)
				return false
			}
		}
		return isFavorite
	}
}
