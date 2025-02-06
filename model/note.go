package model

import "time"

type Note struct {
	Nid              string    `json:"nid" gorm:"column:nid"`
	Uid              uint      `json:"uid" gorm:"column:uid"`
	Cover            string    `json:"cover" gorm:"column:cover"`
	CoverHeight      int       `json:"cover_height" gorm:"column:cover_height"`
	Pics             string    `json:"pics" gorm:"column:pics"`
	Title            string    `json:"title" gorm:"column:title"`
	Content          string    `json:"content" gorm:"column:content"`
	CreatedAt        time.Time `json:"createdAt" gorm:"column:created_at"`
	UpdatedAt        time.Time `json:"updatedAt" gorm:"column:updated_at"`
	Public           uint      `json:"public" gorm:"column:public"`
	CategoryId       uint      `json:"categoryId" gorm:"column:category_id"`
	Tags             string    `json:"tags" gorm:"column:tags"`
	LikesCount       uint      `json:"likesCount" gorm:"column:likes_count"`
	CommentsCount    uint      `json:"commentsCount" gorm:"column:comments_count"`
	CollectionsCount uint      `json:"collectionsCount" gorm:"column:collections_count"`
	SharesCount      uint      `json:"sharesCount" gorm:"column:shares_count"`
	ViewsCount       uint      `json:"viewsCount" gorm:"column:views_count"`
	Status           uint      `json:"status" gorm:"column:status"`
}

func (Note) TableName() string {
	return "notes"
}
