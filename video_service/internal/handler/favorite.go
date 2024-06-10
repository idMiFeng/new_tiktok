package handler

import (
	"api/pb/video"
	"context"
	"fmt"
	"strconv"
	"video_service/global"
	"video_service/internal/model"
	"video_service/internal/service"
)

func (*VideoService) FavoriteAction(ctx context.Context, req *video.FavoriteActionRequest) (resp *video.FavoriteActionResponse, err error) {
	resp = new(video.FavoriteActionResponse)
	if req.UserId == 0 || req.VideoId == 0 || (req.ActionType != 1 && req.ActionType != 2) {
		resp.StatusCode = global.RequestParamErrCode
		resp.StatusMsg = global.GetErrorMessage(global.RequestParamErrCode)
		return resp, nil
	}
	key := fmt.Sprintf("video:user_fav:%s", strconv.FormatInt(req.UserId, 10))
	var redisErr error
	if req.ActionType == 1 {
		// 点赞操作，将视频ID添加到set中
		redisErr = global.Redis.SAdd(ctx, key, req.VideoId).Err()
	} else if req.ActionType == 2 {
		// 取消点赞操作，从set中移除视频ID
		redisErr = global.Redis.SRem(ctx, key, req.VideoId).Err()
	}

	// 检查 Redis 操作结果
	if redisErr != nil {
		resp.StatusCode = global.CacheErrorCode
		resp.StatusMsg = global.GetErrorMessage(global.CacheErrorCode)
		return resp, redisErr
	}
	// 操作成功
	resp.StatusCode = global.SuccessCode
	resp.StatusMsg = global.GetErrorMessage(global.SuccessCode)
	return resp, nil

}

func (*VideoService) FavoriteList(ctx context.Context, req *video.FavoriteListRequest) (resp *video.FavoriteListResponse, err error) {
	resp = new(video.FavoriteListResponse)
	// 检查请求参数
	if req.UserId == 0 {
		resp.StatusCode = global.RequestParamErrCode
		resp.StatusMsg = global.GetErrorMessage(global.RequestParamErrCode)
		return resp, nil
	}
	key := fmt.Sprintf("video:user_fav:%s", strconv.FormatInt(req.UserId, 10))
	// 尝试从 Redis 中获取视频ID集合
	videoIds, err := global.Redis.SMembers(ctx, key).Result()
	if err != nil || len(videoIds) == 0 {
		// 如果 Redis 中没有数据，则从数据库中获取
		videoIdList, dbErr := model.GetFavoriteInstance().GetFavoriteVideoIdsByUser(uint(req.UserId))
		if dbErr != nil {
			resp.StatusCode = global.MYSQLQueryErrCode
			resp.StatusMsg = global.GetErrorMessage(global.MYSQLQueryErrCode)
			return resp, dbErr
		}

		// 将获取到的视频ID存入 Redis
		for _, videoId := range videoIdList {
			global.Redis.SAdd(ctx, key, videoId)
		}

		// 使用数据库中获取到的视频ID列表
		videoIds = make([]string, len(videoIdList))
		for i, videoId := range videoIdList {
			videoIds[i] = strconv.FormatInt(int64(videoId), 10)
		}
	}
	// 获取视频详细信息
	videos := make([]*model.Video, 0, len(videoIds))
	for _, videoIdStr := range videoIds {
		videoId, _ := strconv.ParseInt(videoIdStr, 10, 64)
		v, videoErr := model.GetVideoInstance().GetVideoInfoByVideoId(videoId)
		if videoErr != nil {
			continue // 跳过获取失败的视频
		}
		videos = append(videos, v)
	}
	VideoList := service.BuildVideo(videos, req.UserId)
	// 操作成功
	resp.VideoList = VideoList
	resp.StatusCode = global.SuccessCode
	resp.StatusMsg = global.GetErrorMessage(global.SuccessCode)
	return resp, nil
}
