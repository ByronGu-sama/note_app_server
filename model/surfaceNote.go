package model

type SurfaceNote struct {
	Nid         string `json:"nid" gorm:"column:nid"`
	Uid         uint   `json:"uid" gorm:"column:uid"`
	AvatarUrl   string `json:"avatarUrl" gorm:"column:avatarUrl"`
	Cover       string `json:"cover" gorm:"column:cover"`
	CoverHeight int    `json:"cover_height" gorm:"column:cover_height"`
	Title       string `json:"title" gorm:"column:title"`
	Public      uint   `json:"public" gorm:"column:public"`
	CategoryId  uint   `json:"categoryId" gorm:"column:category_id"`
	Tags        string `json:"tags" gorm:"column:tags"`
	LikesCount  uint   `json:"likesCount" gorm:"column:likes_count"`
}

func (SurfaceNote) TableName() string {
	return "notes"
}
