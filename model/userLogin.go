package model

import "time"

type UserLogin struct {
	Uid                int       `json:"uid" gorm:"column:uid; default:null"`
	Email              string    `json:"email" gorm:"column:email; default:null"`
	Phone              string    `json:"phone"  gorm:"column:phone" binding:"required"`
	Password           string    `json:"password"  gorm:"column:password" binding:"required"`
	LoginFailedTimes   int       `json:"loginFailedTimes"  gorm:"column:loginFailedTimes; default:null"`
	LastLoginFailedAt  time.Time `json:"lastLoginFailedAt" gorm:"column:LastLoginFailedAt; default:null"`
	LastLoginSuccessAt time.Time `json:"lastLoginSuccessAt" gorm:"column:lastLoginSuccessAt; default:null"`
	AccountStatus      int       `json:"accountStatus" gorm:"column:accountStatus; default:null"`
}

func (UserLogin) TableName() string {
	return "user_login"
}
