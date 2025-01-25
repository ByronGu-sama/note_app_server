package repository

import (
	"note_app_server1/global"
	"note_app_server1/model"
	"time"
)

// GetUserInfo 获取用户基本信息
func GetUserInfo(uid uint) (*model.UserInfo, error) {
	var user model.UserInfo
	if err := global.Db.Where("uid = ?", uid).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// GetUserCreationInfo 获取用户创作者信息
func GetUserCreationInfo(uid uint) (*model.UserCreationInfo, error) {
	var info model.UserCreationInfo
	if err := global.Db.Where("uid = ?", uid).First(&info).Error; err != nil {
		return nil, err
	} else {
		return &info, nil
	}
}

// GetUserLoginInfoByPhone 通过手机号获取用户信息
func GetUserLoginInfoByPhone(phone string) (*model.UserLogin, error) {
	var existedUser *model.UserLogin
	if err := global.Db.Where("phone = ?", phone).First(&existedUser).Error; err != nil {
		return nil, err
	}
	return existedUser, nil
}

// GetUserLoginInfoByEmail 通过邮箱获取用户信息
func GetUserLoginInfoByEmail(email string) (*model.UserLogin, error) {
	var existedUser *model.UserLogin
	if err := global.Db.Where("email = ?", email).First(&existedUser).Error; err != nil {
		return nil, err
	}
	return existedUser, nil
}

// UpdateLoginFailedAt 记录上次登陆失败的时间
func UpdateLoginFailedAt(uid uint) {
	global.Db.Model(&model.UserLogin{}).Where("uid = ?", uid).Update("lastLoginFailedAt", time.Now())
}

// UpdateLoginSuccessAt 记录上次登陆成功的时间
func UpdateLoginSuccessAt(uid uint) {
	global.Db.Model(&model.UserLogin{}).Where("uid = ?", uid).Update("lastLoginSuccessAt", time.Now())
}
