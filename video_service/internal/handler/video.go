package handler

import (
	"api/pb/video"
	"context"
	"encoding/json"
	"fmt"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"os"
	"path/filepath"
	"strconv"
	"time"
	"video_service/config"
	"video_service/global"
	"video_service/internal/dao"
	"video_service/internal/dao/mq"
	"video_service/internal/model"
	"video_service/internal/service"
	"video_service/logger"
)

type VideoService struct {
	video.UnimplementedVideoServiceServer // 版本兼容问题
}

func NewVideoService() *VideoService {
	return &VideoService{}
}

func (*VideoService) Feed(ctx context.Context, req *video.FeedRequest) (resp *video.FeedResponse, err error) {
	resp = new(video.FeedResponse)
	var timePoint int64
	if req.LatestTime == -1 {
		timePoint = time.Now().Unix()
	} else {
		timePoint = req.LatestTime
	}
	// 从 Redis 的 ZSET 中获取视频 ID 列表
	zsetKey := "video:publish"
	videoIDs, err := global.Redis.ZRevRangeByScore(ctx, zsetKey, &redis.ZRangeBy{
		Min:    "0",
		Max:    fmt.Sprintf("%d", timePoint),
		Offset: 0,
		Count:  30, // 假设一次取 30 条视频
	}).Result()
	if err != nil {
		logger.Log.Error("Failed to get video IDs from Redis ZSET:", err)
		resp.StatusCode = global.VideoUnExistCode
		resp.StatusMsg = global.GetErrorMessage(global.VideoUnExistCode)
		return resp, err
	}
	var videos []*model.Video
	// 遍历视频 ID 列表，从 Redis 的 HASH 中获取视频信息
	for _, videoID := range videoIDs {
		videoKey := fmt.Sprintf("video:info:%s", videoID)
		videoData, err := global.Redis.HGetAll(ctx, videoKey).Result()
		if err != nil {
			logger.Log.Error("Failed to get video info from Redis HASH:", err)
			continue // 如果获取视频信息失败，跳过该视频
		}
		if len(videoData) == 0 {
			// 如果 Redis 中没有视频信息，从数据库获取并更新到 Redis
			videoIDInt, _ := strconv.Atoi(videoID)
			v, err := model.GetVideoInstance().GetVideoInfoByVideoId(int64(videoIDInt))
			if err != nil {
				logger.Log.Error("Failed to get v from database:", err)
				continue // 如果数据库中也没有视频信息，跳过该视频
			}
			// 更新 Redis 缓存
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
				logger.Log.Error("Failed to update video info in Redis:", err)
			}

			videos = append(videos, v)
		} else {
			// 如果 Redis 中存在视频信息，则直接构建 Video 对象
			videoIDInt, _ := strconv.Atoi(videoData["id"])
			userIDInt, _ := strconv.Atoi(videoData["user_id"])
			favoriteCount, _ := strconv.ParseInt(videoData["favorite_count"], 10, 64)
			commentCount, _ := strconv.ParseInt(videoData["comment_count"], 10, 64)
			v := &model.Video{
				BaseModel: model.BaseModel{
					ID:        uint(videoIDInt),
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				},
				UserId:        uint(userIDInt),
				Title:         videoData["title"],
				PlayURL:       videoData["play_url"],
				CoverURL:      videoData["cover_url"],
				FavoriteCount: favoriteCount,
				CommentCount:  commentCount,
			}
			videos = append(videos, v)
		}

	}

	if req.UserId == -1 {
		// 用户没有登录,返回的视频全部设置为未点赞
		resp.VideoList = service.BuildVideoForFavorite(videos, false)
	} else {
		resp.VideoList = service.BuildVideo(videos, req.UserId)
	}
	// 获取列表中最早发布视频的时间作为下一次请求的时间
	resp.NextTime = videos[len(videos)-1].CreatedAt.Unix()

	resp.StatusCode = global.SuccessCode
	resp.StatusMsg = global.GetErrorMessage(global.SuccessCode)
	return resp, nil
}

func (*VideoService) PublishAction(ctx context.Context, req *video.PublishActionRequest) (resp *video.PublishActionResponse, err error) {
	resp = new(video.PublishActionResponse)
	videoUUID := uuid.New().String()
	videoFileName := fmt.Sprintf("%s.mp4", videoUUID)
	coverFileName := fmt.Sprintf("%s.jpg", videoUUID)
	videoFilePath := filepath.Join("./static", videoFileName)
	coverFilePath := filepath.Join("./static", coverFileName)
	err = os.WriteFile(videoFilePath, req.Data, 0644)
	if err != nil {
		logger.Log.Error("Failed to save video data to temporary file:", err)
		resp.StatusCode = global.VideoUploadErrCode
		resp.StatusMsg = global.GetErrorMessage(global.VideoUploadErrCode)
		return resp, nil
	}
	msg := &model.Video{
		UserId:        uint(req.UserId),
		Title:         req.Title,
		PlayURL:       videoFilePath,
		CoverURL:      coverFilePath,
		FavoriteCount: 0,
		CommentCount:  0,
	}
	msgBody, _ := json.Marshal(msg)

	// 创建 RocketMQ 消息
	rocketmqMsg := &primitive.Message{
		Topic: config.Conf.RocketMQConfig.Topic.VideoPush,
		Body:  msgBody,
	}
	_, err = mq.Producer.SendSync(ctx, rocketmqMsg)
	if err != nil {
		logger.Log.Error("Failed to send message to RocketMQ:", err)
		resp.StatusCode = global.MQSendErrCode
		resp.StatusMsg = global.GetErrorMessage(global.MQSendErrCode)
		return resp, nil
	}
	resp.StatusCode = global.SuccessCode
	resp.StatusMsg = global.GetErrorMessage(global.SuccessCode)
	return resp, nil
}

func (*VideoService) PublishList(ctx context.Context, req *video.PublishListRequest) (resp *video.PublishListResponse, err error) {
	resp = new(video.PublishListResponse)
	key := fmt.Sprintf("%s:%s:%s", "video", "user_work", strconv.FormatInt(req.UserId, 10))

	// 根据用户ID查询视频ID集合
	videoIDs, err := global.Redis.SMembers(dao.Ctx, key).Result()
	if err != nil {
		return nil, fmt.Errorf("缓存错误：%v", err)
	}

	// 查询视频信息
	var videos []*model.Video
	for _, videoID := range videoIDs {
		videoIDInt, err := strconv.Atoi(videoID)
		if err != nil {
			return nil, fmt.Errorf("视频ID转换错误：%v", err)
		}

		v, err := model.GetVideoInstance().GetVideoInfoByVideoId(int64(videoIDInt))
		if err != nil {
			// 处理视频信息查询失败的情况
			continue
		}
		videos = append(videos, v)
	}
	resp.StatusCode = global.SuccessCode
	resp.StatusMsg = global.GetErrorMessage(global.SuccessCode)

	resp.VideoList = service.BuildVideo(videos, req.UserId)

	return resp, nil
}
