package repository

import (
	"errors"
	"gorm.io/gorm"
	"note_app_server/global"
	"note_app_server/model/noteModel"
)

// GetNoteWithNid 获取笔记详情
func GetNoteWithNid(nid string) (*noteModel.NoteDetail, error) {
	var note *noteModel.NoteDetail
	if err := global.Db.Raw("select n.nid as nid, u.uid as uid, u.avatarUrl as avatarUrl, u.username as username, n.pics as pics, n.title as title, n.content as content, n.created_at as created_at, n.updated_at as updated_at, n.public as public, n.category_id as categoryId, n.tags as tags, ni.likes_count as likes_count, ni.comments_count as comments_count, ni.collections_count as collections_count, ni.shares_count as shares_count, ni.views_count as views_count from notes n join user_info u on n.uid = u.uid join notes_info ni on ni.nid = n.nid where n.status = 1 and n.nid = ?", nid).Scan(&note).Error; err != nil {
		return nil, err
	}
	return note, nil
}

// DeleteNoteWithUid 删除笔记
func DeleteNoteWithUid(nid string, uid uint) error {
	result := global.Db.Where("nid = ? and uid = ?", nid, uid).Delete(&noteModel.Note{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("无匹配字段")
	}
	return nil
}

// UpdateNoteWithUid 更新笔记
func UpdateNoteWithUid(n *noteModel.Note) error {
	result := global.Db.Where("nid = ? and uid = ?", n.Nid, n.Uid).Updates(&n)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("update failed")
	}
	return nil
}

// LikeNote 点赞
func LikeNote(nid string, uid uint) error {
	if err := global.Db.Where("uid = ? and nid = ?", uid, nid).First(&noteModel.LikedNotes{}).Error; err == nil {
		return errors.New("has liked")
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	tx := global.Db.Begin()
	if err := tx.Create(&noteModel.LikedNotes{Uid: uid, Nid: nid}).Error; err != nil {
		tx.Rollback()
		return err
	}
	result := tx.Model(&noteModel.NoteInfo{}).Where("nid = ?", nid).UpdateColumn("likes_count", gorm.Expr("likes_count + ?", 1))
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

// CancelLikeNote 取消点赞
func CancelLikeNote(nid string, uid uint) error {
	if err := global.Db.Model(&noteModel.LikedNotes{}).Where("uid = ? and nid = ?", uid, nid).First(&noteModel.LikedNotes{}).Error; err != nil {
		return errors.New("hasn't liked")
	}
	tx := global.Db.Begin()
	result := tx.Model(&noteModel.LikedNotes{}).Where("uid = ? and nid = ?", uid, nid).Delete(&noteModel.LikedNotes{})
	if result.Error != nil {
		tx.Rollback()
		return result.Error
	}
	if result.RowsAffected == 0 {
		tx.Rollback()
		return errors.New("取消点赞失败")
	}
	result = tx.Model(&noteModel.NoteInfo{}).Where("nid = ?", nid).UpdateColumn("likes_count", gorm.Expr("likes_count - ?", 1))
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

// CollectNote 收藏
func CollectNote(nid string, uid uint) error {
	if err := global.Db.Model(&noteModel.CollectedNotes{}).Where("uid = ? and nid = ?", uid, nid).First(&noteModel.CollectedNotes{}).Error; err == nil {
		return errors.New("has collected")
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	tx := global.Db.Begin()
	if err := tx.Create(&noteModel.CollectedNotes{Uid: uid, Nid: nid}).Error; err != nil {
		tx.Rollback()
		return err
	}
	result := tx.Model(&noteModel.NoteInfo{}).Where("nid = ?", nid).UpdateColumn("collections_count", gorm.Expr("collections_count + ?", 1))
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

// CancelCollectNote 取消收藏
func CancelCollectNote(nid string, uid uint) error {
	if err := global.Db.Model(&noteModel.CollectedNotes{}).Where("uid = ? and nid = ?", uid, nid).First(&noteModel.CollectedNotes{}).Error; err != nil {
		return errors.New("hasn't liked")
	}
	tx := global.Db.Begin()
	result := tx.Model(&noteModel.CollectedNotes{}).Where("uid = ? and nid = ?", uid, nid).Delete(&noteModel.CollectedNotes{})
	if result.Error != nil {
		tx.Rollback()
		return result.Error
	}
	if result.RowsAffected == 0 {
		tx.Rollback()
		return errors.New("cancel collect failed")
	}

	result = tx.Model(&noteModel.NoteInfo{}).Where("nid = ?", nid).UpdateColumn("collections_count", gorm.Expr("collections_count - ?", 1))
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

// GetNoteList 查询列表
func GetNoteList(start, limit int) ([]noteModel.SurfaceNote, error) {
	offset := (start - 1) * limit
	var result []noteModel.SurfaceNote
	res := global.Db.Model(&noteModel.SurfaceNote{}).Raw("select n.nid as nid, n.uid as uid, u.username as username, u.avatarUrl as avatarUrl, n.cover as cover, n.cover_height as cover_height, n.title as title, n.public as public, n.category_id as category_id, n.tags as tags, ni.likes_count as like_count from notes n join user_info u on n.uid = u.uid join notes_info ni on n.nid = ni.nid where n.status = 1 limit ?, ?", offset, limit).Scan(&result)
	if res.Error != nil {
		return nil, res.Error
	}
	if res.RowsAffected == 0 {
		return nil, errors.New("没有数据了哦")
	}
	return result, nil
}

// GetNoteListWithUid 查询用户发布的帖子
func GetNoteListWithUid(uid uint, start, limit int) ([]noteModel.SurfaceNote, error) {
	offset := (start - 1) * limit
	var result []noteModel.SurfaceNote
	res := global.Db.Model(&noteModel.SurfaceNote{}).Raw("select n.nid as nid, n.uid as uid, u.username as username,u.avatarUrl as avatarUrl, n.cover as cover, n.cover_height as cover_height, n.title as title, n.public as public, n.category_id as category_id, n.tags as tags, ni.likes_count as like_count from user_info u left join notes n on n.uid = u.uid join notes_info ni on n.uid = u.uid and ni.nid = n.nid where u.uid = ? limit ?, ?", uid, offset, limit).Scan(&result)
	if res.Error != nil {
		return nil, res.Error
	}
	if res.RowsAffected == 0 {
		return nil, errors.New("没有数据了哦")
	}
	return result, nil
}
