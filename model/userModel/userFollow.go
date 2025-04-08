package userModel

type UserFollow struct {
	Uid       int64  `json:"uid" gorm:"column:uid"`
	TargetUid string `json:"target_uid" gorm:"column:target_uid"`
}

func (UserFollow) TableName() string {
	return "user_follow"
}
