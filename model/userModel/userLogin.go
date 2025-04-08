package userModel

import "time"

type UserLogin struct {
	Uid                int64     `json:"uid" gorm:"column:uid; default:null"`
	Email              string    `json:"email" gorm:"column:email; default:null"`
	Phone              string    `json:"phone"  gorm:"column:phone; default:null"`
	Password           string    `json:"password"  gorm:"column:password" binding:"required"`
	LoginFailedTimes   int64     `json:"loginFailedTimes"  gorm:"column:loginFailedTimes; default:null"`
	LastLoginFailedAt  time.Time `json:"lastLoginFailedAt" gorm:"column:lastLoginFailedAt; default:null"`
	LastLoginSuccessAt time.Time `json:"lastLoginSuccessAt" gorm:"column:lastLoginSuccessAt; default:null"`
	AccountStatus      int64     `json:"accountStatus" gorm:"column:accountStatus; default:null"`
}

func (UserLogin) TableName() string {
	return "user_login"
}
