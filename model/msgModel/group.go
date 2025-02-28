package msgModel

import "time"

// Group 群聊结构体
type Group struct {
	GroupId   string    `json:"group_id" gorm:"group_id"`     // 群id
	GroupName string    `json:"group_name" gorm:"group_name"` // 群名
	Desc      string    `json:"desc" gorm:"desc"`             // 简介
	CreatedAt time.Time `json:"created_at" gorm:"created_at"` // 创建时间
	UpdatedAt time.Time `json:"updated_at" gorm:"updated_at"` // 更新时间
	OwnerId   string    `json:"owner_id" gorm:"owner_id"`     //创建者
}

func (Group) TableName() string {
	return "group"
}
