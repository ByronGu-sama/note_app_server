package repository

import (
	"context"
	"gorm.io/gorm"
	"note_app_server/global"
	"note_app_server/model/userModel"
	"note_app_server/service"
	"note_app_server/utils"
	"strconv"
	"time"
)

// GetToken 获取用户相关的信息
func GetToken(ctx context.Context, uid int64) (string, error) {
	// 生成jwt并保存
	token, err := service.GenerateJWT(uid)
	if err != nil {
		return "", err
	}

	err = global.AuthRdb.Set(ctx, strconv.Itoa(int(uid)), token, time.Hour*24*30).Err()
	if err != nil {
		return "", err
	}
	UpdateLoginSuccessAt(ctx, uid)
	return token, nil
}

// GetUserInfo 获取用户基本信息
func GetUserInfo(ctx context.Context, uid int64) (*userModel.UserInfo, error) {
	var user *userModel.UserInfo
	if err := global.Db.WithContext(ctx).Where("uid = ?", uid).First(&user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

// GetUserCreationInfo 获取用户创作者信息
func GetUserCreationInfo(ctx context.Context, uid int64) (*userModel.UserCreationInfo, error) {
	var info *userModel.UserCreationInfo
	if err := global.Db.WithContext(ctx).Where("uid = ?", uid).First(&info).Error; err != nil {
		return nil, err
	} else {
		return info, nil
	}
}

// GetUserLoginInfoByPhone 通过手机号获取用户信息
func GetUserLoginInfoByPhone(ctx context.Context, phone string) (*userModel.UserLogin, error) {
	var existedUser *userModel.UserLogin
	if err := global.Db.WithContext(ctx).Where("phone = ?", phone).First(&existedUser).Error; err != nil {
		return nil, err
	}
	return existedUser, nil
}

// UpdateUserInfo 更新用户信息
func UpdateUserInfo(ctx context.Context, info *userModel.UserInfo) error {
	return utils.WithTx(ctx, global.Db, func(tx *gorm.DB) error {
		var userInfo *userModel.UserInfo
		if info.AvatarUrl == "" {
			if err := tx.Model(userInfo).Where("uid = ?", info.Uid).Updates(map[string]interface{}{"username": info.Username, "birth": info.Birth, "gender": info.Gender, "signature": info.Signature}).Error; err != nil {
				return err
			}
		} else {
			if err := tx.Model(userInfo).Where("uid = ?", info.Uid).Updates(map[string]interface{}{"username": info.Username, "birth": info.Birth, "gender": info.Gender, "signature": info.Signature, "avatarUrl": info.AvatarUrl}).Error; err != nil {
				return err
			}
		}

		if err := tx.Model(userInfo).Where("uid = ?", info.Uid).Update("updateAt", time.Now().Format("2006-01-02 15:04:05")).Error; err != nil {
			return err
		}
		return nil
	})
}

// GetLastAvatarUrl 查询旧头像地址
func GetLastAvatarUrl(ctx context.Context, uid int64) (string, error) {
	var userInfo *userModel.UserInfo
	if err := global.Db.WithContext(ctx).Where("uid = ?", uid).Select("avatarUrl").First(&userInfo).Error; err != nil {
		return "", err
	}
	return userInfo.AvatarUrl, nil
}

// UpdateLoginFailedAt 记录上次登陆失败的时间
func UpdateLoginFailedAt(ctx context.Context, uid int64) {
	global.Db.WithContext(ctx).Model(&userModel.UserLogin{}).Where("uid = ?", uid).Update("lastLoginFailedAt", time.Now())
}

// UpdateLoginSuccessAt 记录上次登陆成功的时间
func UpdateLoginSuccessAt(ctx context.Context, uid int64) {
	global.Db.WithContext(ctx).Model(&userModel.UserLogin{}).Where("uid = ?", uid).Update("lastLoginSuccessAt", time.Now())
}

// GetUserFollowers 获取粉丝列表
func GetUserFollowers(ctx context.Context, uid int64) ([]userModel.FollowUser, error) {
	var userList []userModel.FollowUser
	sql := `select
    ui.uid,
    ui.username,
    ui.avatarUrl
from user_follow uf
left join user_info ui
on uf.uid = ui.uid
where uf.target_uid = ?`
	if err := global.Db.WithContext(ctx).Raw(sql, uid).Scan(&userList).Error; err != nil {
		return nil, err
	}
	return userList, nil
}

// GetUserFollows 获取关注的用户
func GetUserFollows(ctx context.Context, uid int64) ([]userModel.FollowUser, error) {
	var userList []userModel.FollowUser
	sql := `select 
    ui.uid,
    ui.username,
    ui.avatarUrl
from user_follow uf 
left join user_info ui 
on uf.target_uid = ui.uid 
where uf.uid = ?`
	if err := global.Db.WithContext(ctx).Raw(sql, uid).Scan(&userList).Error; err != nil {
		return nil, err
	}
	return userList, nil
}
