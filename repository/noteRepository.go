package repository

import (
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
