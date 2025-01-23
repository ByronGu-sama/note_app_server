package model

import "time"

type UserInfo struct {
	Uid       int       `json:"uid" gorm:"column:uid; default:null"`
	Username  string    `json:"username" gorm:"column:username; default:null"`
	Age       int       `json:"age" gorm:"column:age; default:null"`
	AvatarUrl string    `json:"avatarUrl" gorm:"column:avatarUrl; default:null"`
	Birth     time.Time `json:"birth" gorm:"column:birth; default:null"`
	Gender    int       `json:"gender" gorm:"column:gender; default:null"`
	Address   string    `json:"address" gorm:"column:address; default:null"`
	Language  string    `json:"language" gorm:"column:language; default:null"`
	CreatedAt time.Time `json:"createdAt" gorm:"column:createdAt; default:null"`
	UpdatedAt time.Time `json:"updatedAt" gorm:"column:updatedAt; default:null"`
	UserRole  int       `json:"userRole" gorm:"column:userRole; default:null"`
}

func (UserInfo) TableName() string {
	return "user_info"
}
