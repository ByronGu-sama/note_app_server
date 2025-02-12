package styleModel

type AppStyle struct {
	Uid           uint   `json:"uid" gorm:"uid"`
	ProfileBanner string `json:"profileBanner" gorm:"profile_banner"`
}

func (AppStyle) TableName() string {
	return "app_style"
}
