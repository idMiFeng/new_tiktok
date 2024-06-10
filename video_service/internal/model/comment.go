package model

import (
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
	"sync"
	"video_service/global"
)

// Comment 评论表 /
type Comment struct {
	BaseModel
	VideoId uint   `json:"video_id,string" gorm:"not null;index:comment_video"`
	UserId  uint   `json:"user_id,string" gorm:"not null"`
	Content string `json:"content,string" gorm:"not null"`
}

func (*Comment) TableName() string {
	return "comment"
}

type CommentModel struct {
	DB    *gorm.DB
	Redis *redis.Client
}

var commentModel *CommentModel
var commentOnce sync.Once

// GetCommentInstance 单例实例
func GetCommentInstance() *CommentModel {
	commentOnce.Do(
		func() {
			commentModel = &CommentModel{
				DB:    global.Db,
				Redis: global.Redis,
			}
		})
	return commentModel
}

// CreateComment 创建评论
func (m *CommentModel) CreateComment(commentInfo *Comment) error {
	return m.DB.Create(commentInfo).Error
}

// DeleteComment 删除评论
func (m *CommentModel) DeleteComment(commentId uint) error {
	return m.DB.Where("id = ?", commentId).Unscoped().Delete(&Comment{}).Error
}

// GetCommentAuthorIds 获取指定视频的评论作者ID列表
func (m *CommentModel) GetCommentAuthorIds(videoId uint) ([]uint, error) {
	userIds := make([]uint, 0)
	err := m.DB.Model(&Comment{}).
		Where("video_id = ?", videoId).
		Order("created_at desc").
		Pluck("user_id", &userIds).Error
	return userIds, err
}

// GetCommentList 获取指定视频的评论列表
func (m *CommentModel) GetCommentList(videoId uint) ([]*Comment, error) {
	commentInfos := make([]*Comment, 0)
	err := m.DB.Where("video_id = ?", videoId).
		Order("created_at desc").
		Find(&commentInfos).Error
	return commentInfos, err
}
