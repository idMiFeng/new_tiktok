package model

import (
	"sync"
	"user_service/global"
	"user_service/utils"
)

type User struct {
	BaseModel
	Username        string `json:"username" gorm:"uniqueIndex:username;size:40;not null"` // 设置唯一索引，判断用户名是否重复
	Password        string `json:"password" gorm:"not null"`
	Avatar          string `json:"avatar" gorm:"type:varchar(255);not null"`           // 用户头像
	BackgroundImage string `json:"background_image" gorm:"type:varchar(255);not null"` //背景图片
	Signature       string `json:"signature" gorm:"type:varchar(255);"`                // 个人简介
}

func (*User) TableName() string {
	return "user"
}

type UserController struct{}

var userOnce sync.Once // 单例模式
var userController *UserController

// GetUserInstance 获取单例实例
func GetUserInstance() *UserController {
	userOnce.Do(
		func() {
			userController = &UserController{}
		},
	)
	return userController
}

// CheckUserExist 检查用户是否存在
func (u *UserController) CheckUserExist(username string) (bool, error) {
	var count int64
	err := global.Db.Model(&User{}).Where("username = ?", username).Count(&count).Error
	return count > 0, err
}

// CreateUserInfo 创建用户信息
func (u *UserController) CreateUserInfo(userInfo *User) error {
	return global.Db.Create(userInfo).Error
}

// GetUserInfoByUsername 根据用户名获取用户信息
func (u *UserController) GetUserInfoByUsername(username string) (*User, error) {
	userInfo := &User{}
	err := global.Db.Model(&User{}).Where("username = ?", username).First(userInfo).Error
	return userInfo, err
}

// GetUserInfoByUserID 根据用户ID获取用户信息
func (u *UserController) GetUserInfoByUserID(userID int64) (*User, error) {
	userInfo := &User{}
	err := global.Db.Model(&User{}).Where("user_id = ?", userID).First(userInfo).Error
	return userInfo, err
}

// CheckPassWord 检查密码是否正确
func (u *UserController) CheckPassWord(password string, storePassword string) bool {
	return utils.VerifyPasswordWithHash(password, storePassword)
}

// GetUserList 根据用户ID列表获取用户信息列表
func (u *UserController) GetUserList(userIDs []int64) ([]*User, error) {
	var userInfos []*User
	err := global.Db.Model(&User{}).Where("id IN ?", userIDs).Find(&userInfos).Error
	return userInfos, err
}
