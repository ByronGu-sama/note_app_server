package noteModel

type NoteInfo struct {
	Nid              string `json:"nid" gorm:"column:nid"`
	LikesCount       uint   `json:"likesCount" gorm:"column:likes_count"`
	CommentsCount    uint   `json:"commentsCount" gorm:"column:comments_count"`
	CollectionsCount uint   `json:"collectionsCount" gorm:"column:collections_count"`
	SharesCount      uint   `json:"sharesCount" gorm:"column:shares_count"`
	ViewsCount       uint   `json:"viewsCount" gorm:"column:views_count"`
}

func (NoteInfo) TableName() string {
	return "notes_info"
}
