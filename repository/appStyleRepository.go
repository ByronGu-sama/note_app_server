package repository

import (
	"context"
	"errors"
	"note_app_server/global"
	"note_app_server/model/styleModel"
)

// UpdateProfileBanner 更新用户页banner
func UpdateProfileBanner(ctx context.Context, uid int64, bannerName string) error {
	style := &styleModel.AppStyle{Uid: uid, ProfileBanner: bannerName}
	result := global.Db.WithContext(ctx).Model(&styleModel.AppStyle{}).Where("uid = ?", uid).Updates(map[string]interface{}{"uid": style.Uid, "profile_banner": style.ProfileBanner})
	if result.RowsAffected == 0 {
		return errors.New("更新失败")
	}
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// GetLastBanner 获取用户原始banner图片地址
func GetLastBanner(ctx context.Context, uid int64) (string, error) {
	var style *styleModel.AppStyle
	if err := global.Db.WithContext(ctx).Where("uid = ?", uid).Select("profile_banner").First(&style).Error; err != nil {
		return "", err
	}
	return style.ProfileBanner, nil
}

// GetStyle 获取app风格数据
func GetStyle(ctx context.Context, uid int64) (*styleModel.AppStyle, error) {
	var style *styleModel.AppStyle
	if err := global.Db.WithContext(ctx).Where("uid = ?", uid).First(&style).Error; err != nil {
		return nil, err
	}
	return style, nil
}
