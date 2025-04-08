package commentModel

type CommentsInfo struct {
	Cid        string `json:"cid" gorm:"cid"`
	LikesCount int64  `json:"likesCount" gorm:"likes_count"`
}

func (CommentsInfo) TableName() string {
	return "comments_info"
}
