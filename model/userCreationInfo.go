package model

type UserCreationInfo struct {
	Uid       int `json:"uid" gorm:"column:uid; default:null"`
	Follows   int `json:"follows" gorm:"column:follows; default:null"`
	Followers int `json:"followers"  gorm:"column:followers;default:null"`
	Likes     int `json:"likes" gorm:"column:likes;default:null"`
	Collects  int `json:"collects" gorm:"column:collects;default:null"`
	NoteCount int `json:"noteCount" gorm:"column:noteCount;default:null"`
}

func (UserCreationInfo) TableName() string {
	return "user_creation_info"
}
