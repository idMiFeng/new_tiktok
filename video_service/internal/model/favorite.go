package model

import (
	"gorm.io/gorm"
	"sync"
	"video_service/global"
)

// Favorite 点赞表 /*
type Favorite struct {
	ID      uint `json:"id,string" gorm:"primarykey"`
	UserId  uint `json:"user_id,string" gorm:"uniqueIndex:user_video_id,not null;"`
	VideoId uint `json:"video_id,string" gorm:"uniqueIndex:user_video_id,not null;"`
}

func (*Favorite) TableName() string {
	return "favorite"
}

type FavoriteModel struct {
	DB *gorm.DB
}

var favoriteModel *FavoriteModel
var favoriteOnce sync.Once

func GetFavoriteInstance() *FavoriteModel {
	favoriteOnce.Do(
		func() {
			favoriteModel = &FavoriteModel{
				DB: global.Db,
			}
		})
	return favoriteModel
}

// IsFavoriteVideo 检查用户是否喜欢某个视频
func (m *FavoriteModel) IsFavoriteVideo(userId, videoId uint) (bool, error) {
	var count int64
	err := m.DB.Model(&Favorite{}).Where("user_id = ? AND video_id = ?", userId, videoId).Count(&count).Error
	return count > 0, err
}

// IsFavoriteRecordExist 检查用户对某个视频的喜欢记录是否存在
func (m *FavoriteModel) IsFavoriteRecordExist(userId, videoId uint) (bool, error) {
	var count int64
	err := m.DB.Model(&Favorite{}).Where("user_id = ? AND video_id = ?", userId, videoId).Count(&count).Error
	return count > 0, err
}

// CreateFavoriteRecord 创建喜欢记录
func (m *FavoriteModel) CreateFavoriteRecord(favoriteInfo *Favorite) error {
	return m.DB.Create(favoriteInfo).Error
}

// DeleteFavoriteRecord 删除喜欢记录
func (m *FavoriteModel) DeleteFavoriteRecord(userId, videoId uint) error {
	return m.DB.Where("user_id = ? AND video_id = ?", userId, videoId).Unscoped().Delete(&Favorite{}).Error
}

// GetFavoriteVideoIdsByUser 获取用户喜欢的所有视频ID
func (m *FavoriteModel) GetFavoriteVideoIdsByUser(userId uint) ([]uint, error) {
	var videoIds []uint
	err := m.DB.Model(&Favorite{}).Where("user_id = ?", userId).Pluck("video_id", &videoIds).Error
	return videoIds, err
}

// GetFavoriteCountByUser 获取用户的喜欢数量
func (m *FavoriteModel) GetFavoriteCountByUser(userId uint) (int64, error) {
	var count int64
	err := m.DB.Model(&Favorite{}).Where("user_id = ?", userId).Count(&count).Error
	return count, err
}

// GetFavoriteCountByVideo 获取视频的获赞数量
func (m *FavoriteModel) GetFavoriteCountByVideo(videoId uint) (int64, error) {
	var count int64
	err := m.DB.Model(&Favorite{}).Where("video_id = ?", videoId).Count(&count).Error
	return count, err
}
