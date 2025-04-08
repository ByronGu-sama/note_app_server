package noteModel

// SurfaceNote 预览状态下的note类型
type SurfaceNote struct {
	Nid         string  `json:"nid" gorm:"column:nid"`
	Uid         int64   `json:"uid" gorm:"column:uid"`
	Username    string  `json:"username" gorm:"column:username"`
	AvatarUrl   string  `json:"avatarUrl" gorm:"column:avatarUrl"`
	Cover       string  `json:"cover" gorm:"column:cover"`
	CoverHeight float64 `json:"cover_height" gorm:"column:cover_height"`
	Title       string  `json:"title" gorm:"column:title"`
	Public      int64   `json:"public" gorm:"column:public"`
	CategoryId  int64   `json:"category_id" gorm:"column:category_id"`
	Tags        string  `json:"tags" gorm:"column:tags"`
	LikesCount  int64   `json:"likes_count" gorm:"column:likes_count"`
}

func (SurfaceNote) TableName() string {
	return "notes"
}
