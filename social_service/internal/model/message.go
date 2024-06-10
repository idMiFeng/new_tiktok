package model

import (
	"gorm.io/gorm"
	"social_service/global"
	"sync"
	"time"
)

type Message struct {
	BaseModel
	UserId   uint   `json:"from_user_id" gorm:"index:idx_message;not null"`
	ToUserId uint   `json:"to_user_id" gorm:"index:idx_message;not null"`
	Content  string `json:"content" gorm:"not null"`
}
type MessageModel struct {
	Db *gorm.DB
}

var messageModel *MessageModel
var messageOnce sync.Once // 单例模式

// GetMessageInstance 获取单例实例
func GetMessageInstance() *MessageModel {
	messageOnce.Do(
		func() {
			messageModel = &MessageModel{
				Db: global.Db,
			}
		},
	)
	return messageModel
}

// CreateMessage 创建消息
func (m *MessageModel) CreateMessage(message *Message) error {
	return m.Db.Create(message).Error
}

// GetMessageList 获取按升序排序的消息列表
func (m *MessageModel) GetMessageList(userId, targetId uint, preMsgTime int64) ([]*Message, error) {
	messageList := make([]*Message, 0)
	preMsgTime = preMsgTime / 1000 // Assuming preMsgTime is in milliseconds
	err := m.Db.Model(&Message{}).
		Where("user_id = ? AND to_user_id = ?", userId, targetId).
		Where("created_at > ?", time.Unix(preMsgTime, 0)).
		Order("created_at asc").Find(&messageList).Error
	return messageList, err
}
