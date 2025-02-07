package repository

import (
	"errors"
	"gorm.io/gorm"
	"note_app_server/global"
	"note_app_server/model"
)

// GetNoteWithNid 获取笔记详情
func GetNoteWithNid(nid string) (*model.Note, error) {
	var note *model.Note
	if err := global.Db.Where("nid = ?", nid).First(&note).Error; err != nil {
		return nil, err
	}
	return note, nil
}

// DeleteNoteWithUid 删除笔记
func DeleteNoteWithUid(nid string, uid uint) error {
	result := global.Db.Where("nid = ? and uid = ?", nid, uid).Delete(&model.Note{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("无匹配字段")
	}
	return nil
}

// UpdateNoteWithUid 更新笔记
func UpdateNoteWithUid(note *model.Note) error {
	result := global.Db.Model(&model.Note{}).Where("nid = ? and uid = ?", note.Nid, note.Uid).Updates(&note)
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
	if err := global.Db.Model(&model.LikedNotes{}).Where("uid = ? and nid = ?", uid, nid).First(&model.LikedNotes{}).Error; err == nil {
		return errors.New("has liked")
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	if err := global.Db.Create(&model.LikedNotes{Uid: uid, Nid: nid}).Error; err != nil {
		return err
	}
	return nil
}

// CancelLikeNote 取消点赞
func CancelLikeNote(nid string, uid uint) error {
	if err := global.Db.Model(&model.LikedNotes{}).Where("uid = ? and nid = ?", uid, nid).First(&model.LikedNotes{}).Error; err != nil {
		return errors.New("hasn't liked")
	}
	result := global.Db.Model(&model.LikedNotes{}).Where("uid = ? and nid = ?", uid, nid).Delete(&model.LikedNotes{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("cancel like failed")
	}
	return nil
}

// CollectNote 收藏
func CollectNote(nid string, uid uint) error {
	if err := global.Db.Model(&model.CollectedNotes{}).Where("uid = ? and nid = ?", uid, nid).First(&model.CollectedNotes{}).Error; err == nil {
		return errors.New("has collected")
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	if err := global.Db.Create(&model.CollectedNotes{Uid: uid, Nid: nid}).Error; err != nil {
		return err
	}
	return nil
}

// CancelCollectNote 取消收藏
func CancelCollectNote(nid string, uid uint) error {
	if err := global.Db.Model(&model.CollectedNotes{}).Where("uid = ? and nid = ?", uid, nid).First(&model.CollectedNotes{}).Error; err != nil {
		return errors.New("hasn't liked")
	}
	result := global.Db.Model(&model.CollectedNotes{}).Where("uid = ? and nid = ?", uid, nid).Delete(&model.CollectedNotes{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("cancel collect failed")
	}
	return nil
}

// GetNoteList 查询列表
func GetNoteList(start, limit int) ([]model.SurfaceNote, error) {
	offset := (start - 1) * limit
	var result []model.SurfaceNote
	res := global.Db.Model(&model.SurfaceNote{}).Raw("select n.nid as nid, n.uid as uid, u.avatarUrl as avatarUrl, n.cover as cover, n.cover_height as cover_height, n.title as title, n.public as public, n.category_id as category_id, n.tags as tags, n.likes_count as like_count from notes n join user_info u on n.uid = u.uid where n.status = 1 limit ?, ?", offset, limit).Scan(&result)
	if res.Error != nil {
		return nil, res.Error
	}
	if res.RowsAffected == 0 {
		return nil, errors.New("没有数据了哦")
	}
	return result, nil
}

// GetNoteListWithUid 查询用户发布的帖子
func GetNoteListWithUid(uid uint, start, limit int) ([]model.SurfaceNote, error) {
	offset := (start - 1) * limit
	var result []model.SurfaceNote
	res := global.Db.Model(&model.SurfaceNote{}).Raw("select n.nid as nid, n.uid as uid, u.avatarUrl as avatarUrl, n.cover as cover, n.cover_height as cover_height, n.title as title, n.public as public, n.category_id as category_id, n.tags as tags, n.likes_count as like_count from user_info u left join notes n on n.uid = u.uid where n.status = 1 and u.uid = ? limit ?, ?", uid, offset, limit).Scan(&result)
	if res.Error != nil {
		return nil, res.Error
	}
	if res.RowsAffected == 0 {
		return nil, errors.New("没有数据了哦")
	}
	return result, nil
}
