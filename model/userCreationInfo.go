package model

type UserCreationInfo struct {
	Uid       uint `json:"uid" gorm:"column:uid; default:null"`
	Follows   uint `json:"follows" gorm:"column:follows; default:null"`
	Followers uint `json:"followers"  gorm:"column:followers;default:null"`
	Likes     uint `json:"likes" gorm:"column:likes;default:null"`
	Collects  uint `json:"collects" gorm:"column:collects;default:null"`
	NoteCount uint `json:"noteCount" gorm:"column:noteCount;default:null"`
}

func (UserCreationInfo) TableName() string {
	return "user_creation_info"
}
