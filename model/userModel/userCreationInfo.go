package userModel

type UserCreationInfo struct {
	Uid       int64 `json:"uid" gorm:"column:uid; default:null"`
	Follows   int64 `json:"follows" gorm:"column:follows; default:null"`
	Followers int64 `json:"followers"  gorm:"column:followers;default:null"`
	Likes     int64 `json:"likes" gorm:"column:likes;default:null"`
	Collects  int64 `json:"collects" gorm:"column:collects;default:null"`
	NoteCount int64 `json:"noteCount" gorm:"column:noteCount;default:null"`
}

func (UserCreationInfo) TableName() string {
	return "user_creation_info"
}
