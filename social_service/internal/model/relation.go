package model

import (
	"gorm.io/gorm"
	"social_service/global"
	"sync"
)

type Relation struct {
	BaseModel
	UserId   uint `json:"user_id,string" gorm:"index:idx_relation,not null"`   // 用户ID
	TargetId uint `json:"target_id,string" gorm:"index:idx_relation,not null"` // 目标ID，添加复合索引
	//IsFriend int  `json:"is_friend" gorm:"not null"` // 如果需要保证 relation_id 唯一，可以使用该字段
}

func (*Relation) TableName() string {
	return "relation"
}

type RelationModel struct {
	Db *gorm.DB
}

var relationModel *RelationModel
var relationOnce sync.Once // 单例模式

// GetRelationInstance 获取单例实例
func GetRelationInstance() *RelationModel {
	relationOnce.Do(
		func() {
			relationModel = &RelationModel{
				Db: global.Db,
			}
		},
	)
	return relationModel
}

// CreateRelation 创建Relation记录
func (m *RelationModel) CreateRelation(relation *Relation) error {
	return m.Db.Create(relation).Error
}

// DeleteRelation 删除Relation记录
func (m *RelationModel) DeleteRelation(userId, targetId uint) error {
	return m.Db.Where("user_id = ? AND target_id = ?", userId, targetId).Delete(&Relation{}).Error
}

// GetFollowList 根据用户ID查询关注列表
func (m *RelationModel) GetFollowList(userId uint) ([]int64, error) {
	var follows []int64
	err := m.Db.Model(&Relation{}).Where("user_id = ?", userId).Pluck("target_id", &follows).Error
	return follows, err
}

// GetFollowerList 根据用户ID查询粉丝列表
func (m *RelationModel) GetFollowerList(userId uint) ([]int64, error) {
	var followers []int64
	err := m.Db.Model(&Relation{}).Where("target_id = ?", userId).Pluck("user_id", &followers).Error
	return followers, err
}

// GetFriendList 根据用户ID查询好友列表
func (m *RelationModel) GetFriendList(userId uint) ([]int64, error) {
	var friends []int64
	err := m.Db.Raw(`
		SELECT r1.target_id FROM relation r1
		JOIN relation r2 ON r1.user_id = r2.target_id AND r1.target_id = r2.user_id
		WHERE r1.user_id = ?
	`, userId).Pluck("r1.target_id", &friends).Error
	return friends, err
}

// GetFollowCount 查询用户关注数量
func (m *RelationModel) GetFollowCount(userId uint) (int64, error) {
	var count int64
	err := m.Db.Model(&Relation{}).Where("user_id = ?", userId).Count(&count).Error
	return count, err
}

// GetFollowerCount 查询用户粉丝数量
func (m *RelationModel) GetFollowerCount(targetId uint) (int64, error) {
	var count int64
	err := m.Db.Model(&Relation{}).Where("target_id = ?", targetId).Count(&count).Error
	return count, err
}

// IsFollowing 判断用户是否关注了某人
func (m *RelationModel) IsFollowing(userId, targetId uint) (bool, error) {
	var count int64
	err := m.Db.Model(&Relation{}).Where("user_id = ? AND target_id = ?", userId, targetId).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
