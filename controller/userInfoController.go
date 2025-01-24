package controller

import (
	"note_app_server1/global"
	"note_app_server1/model"
)

// getUserInfo 获取用户基本信息
func getUserInfo(uid uint) (*model.UserInfo, error) {
	var user model.UserInfo
	if err := global.Db.Where("uid = ?", uid).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// getUserCreationInfo 获取用户创作者信息
func getUserCreationInfo(uid uint) (*model.UserCreationInfo, error) {
	var info model.UserCreationInfo
	if err := global.Db.Where("uid = ?", uid).First(&info).Error; err != nil {
		return nil, err
	} else {
		return &info, nil
	}
}
