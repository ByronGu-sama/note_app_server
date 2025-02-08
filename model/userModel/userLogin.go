package userModel

import "time"

type UserLogin struct {
	Uid                uint      `json:"uid" gorm:"column:uid; default:null"`
	Email              string    `json:"email" gorm:"column:email; default:null"`
	Phone              string    `json:"phone"  gorm:"column:phone; default:null"`
	Password           string    `json:"password"  gorm:"column:password" binding:"required"`
	LoginFailedTimes   uint      `json:"loginFailedTimes"  gorm:"column:loginFailedTimes; default:null"`
	LastLoginFailedAt  time.Time `json:"lastLoginFailedAt" gorm:"column:lastLoginFailedAt; default:null"`
	LastLoginSuccessAt time.Time `json:"lastLoginSuccessAt" gorm:"column:lastLoginSuccessAt; default:null"`
	AccountStatus      uint      `json:"accountStatus" gorm:"column:accountStatus; default:null"`
}

func (UserLogin) TableName() string {
	return "user_login"
}
