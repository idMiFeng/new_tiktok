package service

import (
	"api/pb/video"
	"context"
	"fmt"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"strconv"
	"video_service/global"
	"video_service/internal/dao/mq"
	"video_service/internal/model"
)

// BuildComments 构建评论响应
func BuildComments(comments []*model.Comment) []*video.Comment {
	var commentresp []*video.Comment

	for _, comment := range comments {
		commentresp = append(commentresp, &video.Comment{
			Id:         int64(comment.ID),
			UserId:     int64(comment.UserId),
			Content:    comment.Content,
			CreateDate: comment.CreatedAt.Format("01-02"),
		})
	}

	return commentresp
}

// SendCommentDeleteMessage 发送评论删除消息
func SendCommentDeleteMessage(commentID uint) error {
	msg := primitive.NewMessage("comment_delete", []byte(fmt.Sprintf("%d", commentID)))

	// 发送消息
	_, err := mq.Producer.SendSync(context.Background(), msg)
	if err != nil {
		return err
	}

	return nil
}

// CommentMessageHandler 视频消息处理函数
func CommentMessageHandler(ctx context.Context, msgs ...*primitive.MessageExt) (consumer.ConsumeResult, error) {
	for _, msg := range msgs {
		// 处理删除评论消息
		commentIDStr := string(msg.Body)
		commentID, err := strconv.ParseUint(commentIDStr, 10, 64)
		commentKey := fmt.Sprintf("video:comment:%d", commentID)
		if err != nil {
			// 处理错误
			return consumer.ConsumeRetryLater, err
		}
		err = global.Redis.ZRemRangeByScore(ctx, commentKey, fmt.Sprintf("%d", commentID), fmt.Sprintf("%d", commentID)).Err()
		if err != nil {
			return consumer.ConsumeRetryLater, err
		}
	}
	return consumer.ConsumeSuccess, nil
}
