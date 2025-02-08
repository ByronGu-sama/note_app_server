package repository

import (
	"note_app_server/global"
	"note_app_server/model/userModel"
	"time"
)

// GetUserInfo 获取用户基本信息
func GetUserInfo(uid uint) (*userModel.UserInfo, error) {
	var user *userModel.UserInfo
	if err := global.Db.Where("uid = ?", uid).First(&user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

// GetUserCreationInfo 获取用户创作者信息
func GetUserCreationInfo(uid uint) (*userModel.UserCreationInfo, error) {
	var info *userModel.UserCreationInfo
	if err := global.Db.Where("uid = ?", uid).First(&info).Error; err != nil {
		return nil, err
	} else {
		return info, nil
	}
}

// GetUserLoginInfoByPhone 通过手机号获取用户信息
func GetUserLoginInfoByPhone(phone string) (*userModel.UserLogin, error) {
	var existedUser *userModel.UserLogin
	if err := global.Db.Where("phone = ?", phone).First(&existedUser).Error; err != nil {
		return nil, err
	}
	return existedUser, nil
}

// UpdateUserInfo 更新用户信息
func UpdateUserInfo(info *userModel.UserInfo) error {
	var userInfo *userModel.UserInfo
	if err := global.Db.Model(userInfo).Where("uid = ?", info.Uid).Updates(map[string]interface{}{"username": info.Username, "age": info.Age, "birth": info.Birth, "gender": info.Gender, "signature": info.Signature, "address": info.Address}).Error; err != nil {
		return err
	}
	return nil
}

// UpdateUserAvatar 修改头像
func UpdateUserAvatar(uid uint, avatarUrl string) error {
	if err := global.Db.Model(&userModel.UserInfo{}).Where("uid = ?", uid).Updates(map[string]interface{}{"avatarUrl": avatarUrl}).Error; err != nil {
		return err
	}
	return nil
}

// UpdateLoginFailedAt 记录上次登陆失败的时间
func UpdateLoginFailedAt(uid uint) {
	global.Db.Model(&userModel.UserLogin{}).Where("uid = ?", uid).Update("lastLoginFailedAt", time.Now())
}

// UpdateLoginSuccessAt 记录上次登陆成功的时间
func UpdateLoginSuccessAt(uid uint) {
	global.Db.Model(&userModel.UserLogin{}).Where("uid = ?", uid).Update("lastLoginSuccessAt", time.Now())
}
