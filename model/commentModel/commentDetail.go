package commentModel

import "time"

type CommentDetail struct {
	Cid            string          `json:"cid" gorm:"cid"`
	Nid            string          `json:"nid" gorm:"nid"`
	Uid            uint            `json:"uid" gorm:"uid"`
	Username       string          `json:"username" gorm:"username"`
	AvatarUrl      string          `json:"avatarUrl" gorm:"avatar_url"`
	Content        string          `json:"content" gorm:"content"`
	ParentId       string          `json:"parent_id" gorm:"parent_id;default:null"`
	ParentUsername string          `json:"parent_name" gorm:"parent_username"`
	RootId         string          `json:"root_id" gorm:"root_id"`
	CreatedAt      time.Time       `json:"created_at" gorm:"created_at"`
	LikesCount     uint            `json:"likes_count" gorm:"likes_count"`
	Children       []CommentDetail `json:"children" gorm:"-"`
}
