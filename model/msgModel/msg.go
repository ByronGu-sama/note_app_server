package msgModel

import (
	"time"
)

// Message 消息结构体
type Message struct {
	FromId    uint      `json:"from_id" gorm:"from_id"`       // 发送者
	ToId      uint      `json:"to_id" gorm:"to_id"`           //接受者
	Type      int       `json:"type" gorm:"type"`             // 消息类型 群聊 私聊
	Content   string    `json:"content" gorm:"content"`       // 消息内容
	MediaType int       `json:"media_type" gorm:"media_type"` // 媒体类型 文字 图片 视频
	Url       string    `json:"url" gorm:"url"`               // 图片或视频url
	PubTime   time.Time `json:"pub_time" gorm:"pub_time"`     // 发送时间
	GroupId   string    `json:"group_id" gorm:"group_id"`     // 群id
}

func (Message) TableName() string {
	return "message"
}
