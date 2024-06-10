package model

import (
	"gorm.io/gorm"
	"user_service/global"
)

type UserCount struct {
	BaseModel
	UserId         uint  `json:"user_id,string" gorm:"uniqueIndex,not null"`
	User           User  `gorm:"ForeignKey:UserId"`
	FollowCount    int64 `json:"follow_count,string" `    // 关注总数
	FollowerCount  int64 `json:"follower_count,string" `  // 粉丝总数
	TotalFavorited int64 `json:"total_favorited,string" ` // 获赞数量
	WorkCount      int64 `json:"work_count,string" `      // 作品数
	FavoriteCount  int64 `json:"favorite_count,string" `  // 点赞总数
}

func (*UserCount) TableName() string {
	return "user_count"
}

//type UserCountController struct {
//}
//
//var userModel *UserCountController
//var userCountOnce sync.Once // 单例模式
//// GetUserInstance 获取单例实例
//func GetUserCountInstance() *UserCountController {
//	userCountOnce.Do(
//		func() {
//			userModel = &UserCountController{}
//		},
//	)
//	return userModel
//}

func (u *UserController) CreatUserCountInfo(userCountInfo *UserCount) error {
	return global.Db.Create(userCountInfo).Error
}

func (u *UserController) GetUserCountByUserID(userID int64) (*UserCount, error) {
	userCountInfo := &UserCount{}
	err := global.Db.Model(&UserCount{}).Where("user_id = ?", userID).First(userCountInfo).Error
	return userCountInfo, err
}

func (u *UserController) GetUserCountList(userIDs []int64) ([]*UserCount, error) {
	var userCountInfos []*UserCount
	err := global.Db.Model(&UserCount{}).Where("id IN ?", userIDs).Find(&userCountInfos).Error
	return userCountInfos, err
}

// AddFollowCount 增加关注数量
func (u *UserController) AddFollowCount(userID uint) error {
	return global.Db.Model(&UserCount{}).
		Where("user_id = ?", userID).
		Update("follow_count", gorm.Expr("follow_count + ?", 1)).Error
}

// AddFollowerCount 增加粉丝数量
func (u *UserController) AddFollowerCount(userID uint) error {
	return global.Db.Model(&UserCount{}).
		Where("user_id = ?", userID).
		Update("follower_count", gorm.Expr("follower_count + ?", 1)).Error
}

// SubFollowCount 减少关注数量
func (u *UserController) SubFollowCount(userID uint) error {
	return global.Db.Model(&UserCount{}).
		Where("user_id = ?", userID).
		Update("follow_count", gorm.Expr("follow_count - ?", 1)).Error
}

// SubFollowerCount 减少粉丝数量
func (u *UserController) SubFollowerCount(userID uint) error {
	return global.Db.Model(&UserCount{}).
		Where("user_id = ?", userID).
		Update("follower_count", gorm.Expr("follow_count - ?", 1)).Error
}
