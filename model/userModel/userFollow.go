package userModel

type UserFollow struct {
	Uid       uint   `json:"uid" gorm:"column:uid"`
	TargetUid string `json:"target_uid" gorm:"column:target_uid"`
}

func (UserFollow) TableName() string {
	return "user_follow"
}
