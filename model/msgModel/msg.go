package msgModel

import (
	"time"
)

// Message 消息结构体
type Message struct {
	FromId     uint      `json:"from_id"`     // 发送者
	FromName   string    `json:"from_name"`   // 发送者名称
	FromAvatar string    `json:"from_avatar"` // 发送者头像
	ToId       uint      `json:"to_id"`       //接受者
	ToName     string    `json:"to_name"`     // 接受者名称
	Type       int       `json:"type"`        // 消息类型 群聊 私聊
	Content    string    `json:"content"`     // 消息内容
	MediaType  int       `json:"media_type"`  // 媒体类型 文字 图片 视频
	Url        string    `json:"url"`         // 图片或视频url
	PubTime    time.Time `json:"pub_time"`    // 发送时间
	GroupId    string    `json:"group_id"`    // 群id
	Read       bool      `json:"Read"`        // 私聊状态下知否已读
}
