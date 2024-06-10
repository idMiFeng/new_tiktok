package model

import (
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
	"sync"
	"user_service/global"
)

// Video 视频表 /*
type Video struct {
	BaseModel
	UserId        uint   `json:"user_id,string" gorm:"index:user_id,not null"`
	Title         string `json:"title" gorm:"type:varchar(255);not null"`     // 视频标题
	PlayURL       string `json:"play_url" gorm:"type:varchar(255);not null"`  // 视频播放地址
	CoverURL      string `json:"cover_url" gorm:"type:varchar(255);not null"` // 视频封面地址
	FavoriteCount int64  `json:"favorite_count,string" `                      // 视频的点赞总数
	CommentCount  int64  `json:"comment_count,string"`                        // 视频的评论总数
}

func (*Video) TableName() string {
	return "video"
}

type VideoModel struct {
	DB    *gorm.DB
	Redis *redis.Client
}

var videoMedel *VideoModel
var videoOnce sync.Once // 单例模式

// GetVideoInstance 获取单例的实例
func GetVideoInstance() *VideoModel {
	videoOnce.Do(
		func() {
			videoMedel = &VideoModel{
				DB:    global.Db,
				Redis: global.Redis,
			}
		})
	return videoMedel
}

func (v *VideoModel) CreateVideo(videoInfo *Video) error {
	return v.DB.Create(videoInfo).Error
}

func (v *VideoModel) GetVideosByTime(LatestTime int64) ([]*Video, error) {
	videos := make([]*Video, 0)
	err := v.DB.Model(&Video{}).Where("created_at < ?", LatestTime).Order("created_at DESC").Limit(30).Find(&videos).Error
	return videos, err
}

func (v *VideoModel) GetVideosByUserId(userId int64) ([]*Video, error) {
	videos := make([]*Video, 0)
	err := v.DB.Model(&Video{}).Where("user_id = ?", userId).Find(&videos).Error
	return videos, err
}

func (v *VideoModel) GetVideoIdsByUserId(userId int64) ([]int64, error) {
	var videoIds []int64
	err := v.DB.Model(&Video{}).Where("user_id = ?", userId).Pluck("id", &videoIds).Error
	return videoIds, err
}

func (v *VideoModel) GetVideoInfoByVideoId(videoId int64) (*Video, error) {
	var videoInfo *Video
	err := v.DB.Where("id = ?", videoId).First(&videoInfo).Error
	return videoInfo, err
}

func (v *VideoModel) IsVideoExist(videoId int64) (bool, error) {
	var count int64
	err := v.DB.Model(&Video{}).Where("id = ?", videoId).Count(&count).Error
	return count > 0, err
}

func (v *VideoModel) AddVideoFavoriteCount(videoId int64) error {
	return v.DB.Model(&Video{}).Where("id = ?", videoId).Update("favorite_count", gorm.Expr("favorite_count + ?", 1)).Error
}

func (v *VideoModel) SubVideoFavoriteCount(videoId int64) error {
	return v.DB.Model(&Video{}).Where("id = ?", videoId).Update("favorite_count", gorm.Expr("favorite_count - ?", 1)).Error
}
