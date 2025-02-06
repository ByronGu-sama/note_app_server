package model

import "time"

type Comment struct {
	Cid       string    `json:"cid" gorm:"cid"`
	Nid       string    `json:"nid" gorm:"nid"`
	Uid       uint      `json:"uid" gorm:"uid"`
	Content   string    `json:"content" gorm:"content"`
	ParentId  string    `json:"parent_id" gorm:"parent_id"`
	CreatedAt time.Time `json:"created_at" gorm:"created_at"`
}

func (Comment) TableName() string {
	return "comments"
}
