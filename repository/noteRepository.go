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
func GetNoteList(start int, limit int) ([]model.Note, error) {
	offset := (start - 1) * limit
	var result []model.Note
	res := global.Db.Model(&model.Note{}).Offset(offset).Limit(limit).Find(&result)
	if res.Error != nil {
		return nil, res.Error
	}
	if res.RowsAffected == 0 {
		return nil, errors.New("没有数据了哦")
	}
	return result, nil
}
