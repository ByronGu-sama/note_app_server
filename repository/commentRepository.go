package repository

import (
	"errors"
	"gorm.io/gorm"
	"note_app_server/global"
	"note_app_server/model/commentModel"
	"note_app_server/model/noteModel"
	"note_app_server/utils"
)

// NewComment 创建评论
func NewComment(cmt *commentModel.Comment, cmtInfo *commentModel.CommentsInfo) (*commentModel.CommentDetail, error) {
	tx := global.Db.Begin()
	if err := tx.Model(&commentModel.Comment{}).Create(cmt).Error; err != nil {
		tx.Rollback()
		return nil, err
	}
	if err := tx.Model(&commentModel.CommentsInfo{}).Create(cmtInfo).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	result := tx.Model(&noteModel.NoteInfo{}).Where("nid = ?", cmt.Nid).UpdateColumn("comments_count", gorm.Expr("comments_count + ?", 1))
	if result.Error != nil {
		tx.Rollback()
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		tx.Rollback()
		return nil, errors.New("更新数据失败")
	}
	tx.Commit()
	var newComment *commentModel.CommentDetail
	if err := global.Db.Raw(`SELECT 
    c.cid AS cid,
    c.nid AS nid,
    c.uid AS uid,
    ui.username AS username,
    ui.avatarUrl AS avatar_url,
    c.content AS content,
    c.root_id AS root_id,
    c.created_at AS created_at,
    ci.likes_count AS likes_count
FROM user_info ui 
    left join comments c on ui.uid = c.uid 
    left join comments_info ci on c.cid = ci.cid
where c.cid = ?`, cmt.Cid).First(&newComment).Error; err != nil {
		return nil, err
	}
	newComment.AvatarUrl = utils.AddAvatarPrefix(newComment.AvatarUrl)
	return newComment, nil
}

// DeleteComment 删除评论
func DeleteComment(uid uint, cid string) error {
	cmt := new(commentModel.Comment)
	if err := global.Db.Where("cid = ? and uid = ?", cid, uid).First(&cmt).Error; err != nil {
		return err
	}
	tx := global.Db.Begin()
	result := tx.Model(&commentModel.Comment{}).Where("cid = ? and uid = ?", cid, uid).Delete(&commentModel.Comment{})
	if result.RowsAffected == 0 {
		tx.Rollback()
		return errors.New("无匹配记录")
	}
	if result.Error != nil {
		tx.Rollback()
		return result.Error
	}

	result = tx.Model(&noteModel.NoteInfo{}).Where("nid = ?", cmt.Nid).UpdateColumn("comments_count", gorm.Expr("comments_count - ?", 1))
	if result.Error != nil {
		tx.Rollback()
		return result.Error
	}
	if result.RowsAffected == 0 {
		tx.Rollback()
		return errors.New("更新数据失败")
	}
	tx.Commit()
	return nil
}

// LikeComment 点赞评论
func LikeComment(uid uint, cid string) error {
	if err := global.Db.Where("uid = ? and cid = ?", uid, cid).First(&commentModel.LikedComment{}).Error; err == nil {
		return errors.New("已点赞")
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	tx := global.Db.Begin()
	result := tx.Create(&commentModel.LikedComment{Uid: uid, Cid: cid})
	if result.Error != nil {
		tx.Rollback()
		return result.Error
	}

	result = tx.Model(&commentModel.CommentsInfo{}).Where("cid = ?", cid).UpdateColumn("likes_count", gorm.Expr("likes_count + ?", 1))
	if result.Error != nil {
		tx.Rollback()
		return result.Error
	}
	if result.RowsAffected == 0 {
		tx.Rollback()
		return errors.New("更新数据失败")
	}
	tx.Commit()
	return nil
}

// DislikeComment 取消点赞评论
func DislikeComment(uid uint, cid string) error {
	if err := global.Db.Where("uid = ? and cid = ?", uid, cid).First(&commentModel.LikedComment{}).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		return errors.New("未点赞")
	} else if err != nil {
		return err
	}

	tx := global.Db.Begin()
	result := tx.Where("uid = ? and cid = ?", uid, cid).Delete(&commentModel.LikedComment{})
	if result.Error != nil {
		tx.Rollback()
		return result.Error
	}
	if result.RowsAffected == 0 {
		tx.Rollback()
		return errors.New("取消点赞失败")
	}

	result = tx.Model(&commentModel.CommentsInfo{}).Where("cid = ?", cid).UpdateColumn("likes_count", gorm.Expr("likes_count - ?", 1))
	if result.Error != nil {
		tx.Rollback()
		return result.Error
	}
	if result.RowsAffected == 0 {
		tx.Rollback()
		return errors.New("更新数据失败")
	}
	tx.Commit()
	return nil
}

// GetNoteCommentsList 获取笔记评论列表
func GetNoteCommentsList(nid string, page, limit int) ([]commentModel.CommentDetail, error) {
	offset := (page - 1) * limit
	var commentsList []commentModel.CommentDetail
	if err := global.Db.Raw(`
		WITH RECURSIVE CommentTree AS(
		(SELECT * FROM comments WHERE parent_id IS NULL LIMIT ?, ?)
		UNION ALL
		SELECT c.* FROM comments c 
		JOIN CommentTree ct ON ct.cid = c.parent_id) 
		SELECT 
			ctResult.cid AS cid, 
			ctResult.nid AS nid, 
			ctResult.uid AS uid, 
			ui.username AS username, 
			ui.avatarUrl AS avatar_url, 
			ctResult.content AS content, 
			ctResult.parent_id AS parent_id, 
			ctResult.root_id AS root_id,
			parentUser.username AS parent_username, 
			ctResult.created_at AS created_at, 
			ci.likes_count AS likes_count 
		FROM CommentTree ctResult 
		JOIN comments_info ci ON ctResult.cid = ci.cid 
		JOIN user_info ui ON ui.uid = ctResult.uid 
		LEFT JOIN comments parentComment ON parentComment.cid = ctResult.parent_id 
		LEFT JOIN user_info parentUser ON parentUser.uid = parentComment.uid
		WHERE ctResult.nid = ?`,
		offset, limit, nid).Scan(&commentsList).Error; err != nil {
		return nil, err
	}
	return commentsList, nil
}
