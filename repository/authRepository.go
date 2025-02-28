package repository

import (
	"note_app_server/global"
	"note_app_server/model/styleModel"
	"note_app_server/model/userModel"
)

func RegisterUser(user *userModel.UserLogin) error {
	tx := global.Db.Begin()
	if err := tx.Create(&user).Error; err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()

	tx = global.Db.Begin()
	newUser, err := GetUserLoginInfoByPhone(user.Phone) //获取系统生成的用户uid
	if err != nil {
		tx.Rollback()
		return err
	}

	userInfo := &userModel.UserInfo{Uid: newUser.Uid}
	userCreationInfo := &userModel.UserCreationInfo{Uid: newUser.Uid}
	style := &styleModel.AppStyle{Uid: newUser.Uid}

	// 创建用户详细信息
	if err = tx.Create(&userInfo).Error; err != nil {
		tx.Rollback()
		return err
	}
	// 创建用户创造者信息
	if err = tx.Create(&userCreationInfo).Error; err != nil {
		tx.Rollback()
		return err
	}
	// 创建默认style
	if err = tx.Create(&style).Error; err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}
