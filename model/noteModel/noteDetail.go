package noteModel

import "time"

// NoteDetail note详情，包含作者信息
type NoteDetail struct {
	Nid              string    `json:"nid" gorm:"column:nid"`
	Uid              int64     `json:"uid" gorm:"column:uid"`
	AvatarUrl        string    `json:"avatarUrl" gorm:"column:avatarUrl"`
	Username         string    `json:"username" gorm:"column:username"`
	Pics             string    `json:"pics" gorm:"column:pics"`
	Title            string    `json:"title" gorm:"column:title"`
	Content          string    `json:"content" gorm:"column:content"`
	CreatedAt        time.Time `json:"createdAt" gorm:"column:created_at"`
	UpdatedAt        time.Time `json:"updatedAt" gorm:"column:updated_at"`
	Public           int64     `json:"public" gorm:"column:public"`
	CategoryId       int64     `json:"categoryId" gorm:"column:category_id"`
	Tags             string    `json:"tags" gorm:"column:tags"`
	LikesCount       int64     `json:"likesCount" gorm:"column:likes_count"`
	CommentsCount    int64     `json:"commentsCount" gorm:"column:comments_count"`
	CollectionsCount int64     `json:"collectionsCount" gorm:"column:collections_count"`
	SharesCount      int64     `json:"sharesCount" gorm:"column:shares_count"`
	ViewsCount       int64     `json:"viewsCount" gorm:"column:views_count"`
}
