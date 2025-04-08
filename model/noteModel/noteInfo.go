package noteModel

type NoteInfo struct {
	Nid              string `json:"nid" gorm:"column:nid"`
	LikesCount       int64  `json:"likesCount" gorm:"column:likes_count"`
	CommentsCount    int64  `json:"commentsCount" gorm:"column:comments_count"`
	CollectionsCount int64  `json:"collectionsCount" gorm:"column:collections_count"`
	SharesCount      int64  `json:"sharesCount" gorm:"column:shares_count"`
	ViewsCount       int64  `json:"viewsCount" gorm:"column:views_count"`
}

func (NoteInfo) TableName() string {
	return "notes_info"
}
