package service

import (
	"api/pb/video"
	"context"
	"encoding/json"
	"fmt"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/go-redis/redis/v8"
	"os/exec"
	"strconv"
	"video_service/global"
	"video_service/internal/model"
	"video_service/logger"
)

func BuildVideoForFavorite(videos []*model.Video, isFavorite bool) []*video.Video {
	var videoResp []*video.Video

	for _, v := range videos {

		videoResp = append(videoResp, &video.Video{
			Id:            int64(v.ID),
			AuthId:        int64(v.UserId),
			PlayUrl:       v.PlayURL,
			CoverUrl:      v.CoverURL,
			FavoriteCount: v.FavoriteCount,
			CommentCount:  v.CommentCount,
			IsFavorite:    isFavorite,
			Title:         v.Title,
		})
	}

	return videoResp
}

func BuildVideo(videos []*model.Video, userId int64) []*video.Video {
	var videoResp []*video.Video

	for _, v := range videos {
		// 查询是否有喜欢的缓存，如果有，比对缓存，如果没有，构建缓存再查缓存
		favorite := IsFavorite(uint(userId), v.ID)

		videoResp = append(videoResp, &video.Video{
			Id:            int64(v.ID),
			AuthId:        int64(v.UserId),
			PlayUrl:       v.PlayURL,
			CoverUrl:      v.CoverURL,
			FavoriteCount: v.FavoriteCount,
			CommentCount:  v.CommentCount,
			IsFavorite:    favorite,
			Title:         v.Title,
		})
	}

	return videoResp
}

func VideoMessageHandler(ctx context.Context, msgs ...*primitive.MessageExt) (consumer.ConsumeResult, error) {
	for _, msg := range msgs {
		var v model.Video
		err := json.Unmarshal(msg.Body, &v)
		if err != nil {
			logger.Log.Error("Failed to unmarshal message:", err)
			return consumer.ConsumeRetryLater, nil
		}

		coverFilePath := v.CoverURL
		videoFilePath := v.PlayURL

		// 使用 FFmpeg 截取视频第一帧
		err = extractFirstFrame(videoFilePath, coverFilePath)
		if err != nil {
			logger.Log.Error("Failed to extract first frame:", err)
			return consumer.ConsumeRetryLater, nil
		}

		// 将视频信息保存到数据库
		err = model.GetVideoInstance().CreateVideo(&v)
		if err != nil {
			logger.Log.Error("Failed to save v to database:", err)
			return consumer.ConsumeRetryLater, nil
		}
		// 将视频信息保存到 Redis
		videoKey := fmt.Sprintf("video:info:%s", strconv.Itoa(int(v.ID)))
		videoData := map[string]interface{}{
			"id":             v.ID,
			"user_id":        v.UserId,
			"title":          v.Title,
			"play_url":       v.PlayURL,
			"cover_url":      v.CoverURL,
			"favorite_count": v.FavoriteCount,
			"comment_count":  v.CommentCount,
		}
		err = global.Redis.HMSet(ctx, videoKey, videoData).Err()
		if err != nil {
			logger.Log.Error("Failed to save video to Redis:", err)
			return consumer.ConsumeRetryLater, nil
		}
		// 将视频信息保存到 Redis ZSET
		zsetKey := "video:publish"
		err = global.Redis.ZAdd(ctx, zsetKey, &redis.Z{
			Score:  float64(v.CreatedAt.Unix()),
			Member: v.ID,
		}).Err()
		if err != nil {
			logger.Log.Error("Failed to save video to Redis ZSET:", err)
			return consumer.ConsumeRetryLater, nil
		}
	}

	return consumer.ConsumeSuccess, nil
}

func extractFirstFrame(videoPath, coverPath string) error {
	cmd := exec.Command("ffmpeg", "-i", videoPath, "-vf", "thumbnail", "-frames:v", "1", coverPath)
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("ffmpeg command failed: %v", err)
	}
	return nil
}
