package service

import (
	"api/pb/video"
	"video_service/internal/model"
)

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
