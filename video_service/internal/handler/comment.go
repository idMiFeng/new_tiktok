package handler

import (
	"api/pb/video"
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"strconv"
	"time"
	"video_service/global"
	"video_service/internal/model"
	"video_service/internal/service"
	"video_service/logger"
)

func (*VideoService) CommentAction(ctx context.Context, req *video.CommentActionRequest) (resp *video.CommentActionResponse, err error) {
	resp = new(video.CommentActionResponse)
	commentModel := model.GetCommentInstance()
	commentKey := fmt.Sprintf("video:comment:%d", req.VideoId)
	timeNow := time.Now()

	switch req.ActionType {
	case 1: // 添加评论
		// ① 将评论信息添加到数据库中
		commentId := uuid.New().ID()
		commentInfo := &model.Comment{
			BaseModel: model.BaseModel{ID: uint(commentId), CreatedAt: timeNow, UpdatedAt: timeNow},
			VideoId:   uint(req.VideoId),
			UserId:    uint(req.UserId),
			Content:   req.CommentText,
		}
		err = commentModel.CreateComment(commentInfo)
		if err != nil {
			return &video.CommentActionResponse{StatusCode: global.MYSQLQueryErrCode, StatusMsg: global.GetErrorMessage(global.MYSQLQueryErrCode)}, nil
		}

		// ② 直接对redis中video_info 的 comment_count进行更改
		videoKey := fmt.Sprintf("video:info:%d", req.VideoId)
		err = global.Redis.HIncrBy(ctx, videoKey, "comment_count", 1).Err()
		if err != nil {
			logger.Log.Error("Failed to increment comment count in Redis:", err)
			return &video.CommentActionResponse{StatusCode: global.CacheErrorCode, StatusMsg: global.GetErrorMessage(global.CacheErrorCode)}, nil
		}
		// ③ 将comment加入redis，用zset
		content := fmt.Sprintf("%s:%d:%s", strconv.Itoa(int(commentId)), req.UserId, req.CommentText)
		err = global.Redis.ZAdd(ctx, commentKey, &redis.Z{
			Score:  float64(timeNow.Unix()),
			Member: content,
		}).Err()
		if err != nil {
			return &video.CommentActionResponse{StatusCode: global.CacheErrorCode, StatusMsg: global.GetErrorMessage(global.CacheErrorCode)}, nil
		}

		// 模型转换，格式化时间
		createDate := timeNow.Format("01-02")
		resp = &video.CommentActionResponse{
			StatusCode: global.SuccessCode,
			StatusMsg:  global.GetErrorMessage(global.SuccessCode),
			Comment: &video.Comment{
				Id:         int64(commentId),
				UserId:     req.UserId,
				Content:    req.CommentText,
				CreateDate: createDate,
			},
		}

	case 2: // 删除评论
		// ① 从数据库中删除评论数据
		err = commentModel.DeleteComment(uint(req.CommentId))
		if err != nil {
			return &video.CommentActionResponse{StatusCode: global.MYSQLQueryErrCode, StatusMsg: global.GetErrorMessage(global.MYSQLQueryErrCode)}, nil
		}

		// 删除缓存
		err = global.Redis.ZRemRangeByScore(ctx, commentKey, fmt.Sprintf("%d", req.CommentId), fmt.Sprintf("%d", req.CommentId)).Err()
		if err != nil {
			return &video.CommentActionResponse{StatusCode: global.CacheErrorCode, StatusMsg: global.GetErrorMessage(global.CacheErrorCode)}, nil
		}

		// ② 对视频的comment_count-1
		videoKey := fmt.Sprintf("video:info:%d", req.VideoId)
		err = global.Redis.HIncrBy(ctx, videoKey, "comment_count", -1).Err()
		if err != nil {
			logger.Log.Error("Failed to decrement comment count in Redis:", err)
			return &video.CommentActionResponse{StatusCode: 1, StatusMsg: "Failed to update comment count in Redis"}, nil
		}

		resp = &video.CommentActionResponse{StatusCode: global.SuccessCode, StatusMsg: global.GetErrorMessage(global.SuccessCode)}

	default:
		return &video.CommentActionResponse{StatusCode: global.RequestParamErrCode, StatusMsg: global.GetErrorMessage(global.RequestParamErrCode)}, nil
	}

	return resp, nil
}
func (*VideoService) CommentList(ctx context.Context, req *video.CommentListRequest) (resp *video.CommentListResponse, err error) {
	resp = &video.CommentListResponse{}
	commentKey := fmt.Sprintf("video:comment:%d", req.VideoId)
	var comments []*model.Comment

	// 尝试从 Redis 中获取评论列表
	commentData, err := global.Redis.Get(ctx, commentKey).Result()
	if err == redis.Nil {
		// Redis 未命中缓存，从数据库中获取
		comments, err = model.GetCommentInstance().GetCommentList(uint(req.VideoId))
		if err != nil {
			resp.StatusCode = global.MYSQLQueryErrCode
			resp.StatusMsg = global.GetErrorMessage(global.MYSQLQueryErrCode)
			return resp, err
		}

		// 将评论列表缓存到 Redis 中
		commentDataBytes, err := json.Marshal(comments)
		if err != nil {
			resp.StatusCode = 1
			resp.StatusMsg = "Failed to marshal comments"
			return resp, err
		}

		err = global.Redis.Set(ctx, commentKey, commentDataBytes, 10*time.Minute).Err()
		if err != nil {
			resp.StatusCode = global.RequestParamErrCode
			resp.StatusMsg = global.GetErrorMessage(global.RequestParamErrCode)
			return resp, err
		}
	} else if err != nil {
		resp.StatusCode = global.CacheErrorCode
		resp.StatusMsg = "Failed to get comments from Redis"
		return resp, err
	} else {
		// Redis 命中缓存，反序列化评论列表
		err = json.Unmarshal([]byte(commentData), &comments)
		if err != nil {
			resp.StatusCode = 1
			resp.StatusMsg = "Failed to unmarshal comments from Redis"
			return resp, err
		}
	}

	// 构建响应
	resp.StatusCode = global.SuccessCode
	resp.StatusMsg = global.GetErrorMessage(global.SuccessCode)
	resp.CommentList = service.BuildComments(comments)
	return resp, nil
}
