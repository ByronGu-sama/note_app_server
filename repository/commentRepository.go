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

// GetNoteCommentsList 获取笔记评论列表
func GetNoteCommentsList(nid string, page, limit int) ([]commentModel.CommentDetail, error) {
	offset := (page - 1) * limit
	var commentsList []commentModel.CommentDetail
	if err := global.Db.Raw(`
		SELECT
			c.cid AS cid,
			c.nid AS nid,
			c.uid AS uid,
			ui.username AS username,
			ui.avatarUrl AS avatar_url,
			c.content AS content,
			c.parent_id AS parent_id,
			c.root_id AS root_id,
			c.created_at AS created_at,
			ci.likes_count AS likes_count
		FROM comments c
		JOIN comments_info ci ON c.cid = ci.cid
		JOIN user_info ui ON ui.uid = c.uid
		where c.nid = ? and c.parent_id is null limit ?, ?`,
		nid, offset, limit).Scan(&commentsList).Error; err != nil {
		return nil, err
	}
	return commentsList, nil
}

// GetSubCommentsList 获取子评论
func GetSubCommentsList(nid, rootId string, page, limit int) ([]commentModel.CommentDetail, error) {
	offset := (page - 1) * limit
	var commentsList []commentModel.CommentDetail
	if err := global.Db.Raw(`
		SELECT
			c.cid AS cid,
			c.nid AS nid,
			c.uid AS uid,
			ui.username AS username,
			ui.avatarUrl AS avatar_url,
			c.content AS content,
			c.parent_id AS parent_id,
			parent_ui.username AS parent_username,
			c.root_id AS root_id,
			c.created_at AS created_at,
			ci.likes_count AS likes_count
		FROM comments c
		JOIN comments_info ci ON c.cid = ci.cid
		JOIN user_info ui ON ui.uid = c.uid
		LEFT JOIN comments parent_c ON parent_c.cid = c.parent_id
		LEFT JOIN user_info parent_ui ON parent_ui.uid = parent_c.uid
		where c.nid = ? and c.root_id = ? and c.root_id != c.cid and c.parent_id is not null limit ?, ?
		`,
		nid, rootId, offset, limit).Scan(&commentsList).Error; err != nil {
		return nil, err
	}
	return commentsList, nil
}

// DeleteComment 删除评论
func DeleteComment(nid, cid string, uid uint) error {
	tx := global.Db.Begin()
	if err := tx.Delete(&commentModel.CommentDetail{Cid: cid, Uid: uid}).Error; err != nil {
		tx.Rollback()
		return err
	}
	result := tx.Model(&noteModel.NoteInfo{}).Where("nid = ?", nid).UpdateColumn("comments_count", gorm.Expr("comments_count - ?", 1))
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
