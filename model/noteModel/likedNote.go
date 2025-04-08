package noteModel

type LikedNotes struct {
	Uid int64  `json:"uid" gorm:"uid"`
	Nid string `json:"nid" gorm:"nid"`
}

func (LikedNotes) TableName() string {
	return "liked_note"
}
