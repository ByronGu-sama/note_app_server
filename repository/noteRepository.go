package repository

import (
	"errors"
	"note_app_server/global"
	"note_app_server/model"
)

// GetNoteWithNid 获取笔记详情
func GetNoteWithNid(nid uint) (*model.Note, error) {
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
