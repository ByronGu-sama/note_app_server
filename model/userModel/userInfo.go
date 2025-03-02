package userModel

import "time"

type UserInfo struct {
	Uid       uint      `json:"uid" gorm:"column:uid; default:null"`
	Username  string    `json:"username" gorm:"column:username; default:momo"`
	AvatarUrl string    `json:"avatarUrl" gorm:"column:avatarUrl; default:null"`
	Birth     time.Time `json:"birth" gorm:"column:birth; default:null"`
	Gender    string    `json:"gender" gorm:"column:gender; default:null"`
	Signature string    `json:"signature" gorm:"column:signature; default:null"`
	Address   string    `json:"address" gorm:"column:address; default:null"`
	Language  string    `json:"language" gorm:"column:language; default:null"`
	CreateAt  time.Time `json:"createAt" gorm:"column:createAt; default:null"`
	UpdateAt  time.Time `json:"updateAt" gorm:"column:updateAt; default:null"`
	UserRole  uint      `json:"userRole" gorm:"column:userRole; default:null"`
}

func (UserInfo) TableName() string {
	return "user_info"
}

// FollowUser 用户关注列表&粉丝列表类型
type FollowUser struct {
	Uid       uint   `json:"uid" gorm:"column:uid"`
	Username  string `json:"username" gorm:"column:username"`
	AvatarUrl string `json:"avatarUrl" gorm:"column:avatarUrl"`
}
