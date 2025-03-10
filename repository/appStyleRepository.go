package repository

import (
	"errors"
	"note_app_server/global"
	"note_app_server/model/styleModel"
)

// UpdateProfileBanner 更新用户页banner
func UpdateProfileBanner(uid uint, bannerName string) error {
	style := &styleModel.AppStyle{Uid: uid, ProfileBanner: bannerName}
	result := global.Db.Model(&styleModel.AppStyle{}).Where("uid = ?", uid).Updates(map[string]interface{}{"uid": style.Uid, "profile_banner": style.ProfileBanner})
	if result.RowsAffected == 0 {
		return errors.New("更新失败")
	}
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// GetLastBanner 获取用户原始banner图片地址
func GetLastBanner(uid uint) (string, error) {
	var style *styleModel.AppStyle
	if err := global.Db.Where("uid = ?", uid).Select("profile_banner").First(&style).Error; err != nil {
		return "", err
	}
	return style.ProfileBanner, nil
}

// GetStyle 获取app风格数据
func GetStyle(uid uint) (*styleModel.AppStyle, error) {
	var style *styleModel.AppStyle
	if err := global.Db.Where("uid = ?", uid).First(&style).Error; err != nil {
		return nil, err
	}
	return style, nil
}
