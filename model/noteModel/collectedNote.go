package noteModel

type CollectedNotes struct {
	Uid uint   `json:"uid" gorm:"uid"`
	Nid string `json:"nid" gorm:"nid"`
}

func (CollectedNotes) TableName() string {
	return "collected_notes"
}
