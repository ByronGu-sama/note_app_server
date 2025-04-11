package repository

import (
	"context"
	"gorm.io/gorm"
	"note_app_server/global"
	"note_app_server/model/styleModel"
	"note_app_server/model/userModel"
	"note_app_server/utils"
)

func RegisterUser(ctx context.Context, user *userModel.UserLogin) error {
	return utils.WithTx(ctx, global.Db, func(tx *gorm.DB) error {
		if err := tx.WithContext(ctx).Create(&user).Error; err != nil {
			return err
		}
		userInfo := &userModel.UserInfo{Uid: user.Uid}
		userCreationInfo := &userModel.UserCreationInfo{Uid: user.Uid}
		style := &styleModel.AppStyle{Uid: user.Uid}

		// 创建用户详细信息
		if err := tx.WithContext(ctx).Create(&userInfo).Error; err != nil {
			return err
		}
		// 创建用户创造者信息
		if err := tx.WithContext(ctx).Create(&userCreationInfo).Error; err != nil {
			return err
		}
		// 创建默认style
		if err := tx.WithContext(ctx).Create(&style).Error; err != nil {
			return err
		}
		return nil
	})
}
