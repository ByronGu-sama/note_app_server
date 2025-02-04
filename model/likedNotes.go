package model

type LikedNotes struct {
	Uid uint   `json:"uid" gorm:"uid"`
	Nid string `json:"nid" gorm:"nid"`
}

func (LikedNotes) TableName() string {
	return "liked_notes"
}
