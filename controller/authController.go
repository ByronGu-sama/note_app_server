package controller

import (
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"note_app_server1/global"
	"note_app_server1/model"
)

func Register(ctx *gin.Context) {
	var user model.UserLogin
	// 检查必要字段是否缺失
	if err := ctx.ShouldBind(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "注册失败",
		})
		return
	}

	//检查用户名或其他唯一字段是否已存在
	var existedUser model.UserLogin
	if err := global.Db.Where("phone = ?", user.Phone).First(&existedUser).Error; err == nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "该手机号已注册",
		})
		return
	}

	// 检查密码是否为空
	if user.Password == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "缺少必要信息",
		})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "failed to hash password",
		})
		return
	}
	user.Password = string(hashedPassword)

	// 创建jwt
	//

	// 使用事务创建用户登录信息和详细信息
	tx := global.Db.Begin()
	if err := tx.Create(&user).Error; err != nil {
		tx.Rollback()
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to create user",
		})
		return
	}
	tx.Commit()
	tx = global.Db.Begin()
	if err := tx.Where("phone = ?", user.Phone).First(&user).Error; err != nil {
		tx.Rollback()
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to create user",
		})
		return
	}
	var userInfo model.UserInfo
	userInfo.Uid = user.Uid
	userInfo.Username = "momo"
	if err := tx.Create(&userInfo).Error; err != nil {
		tx.Rollback()
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to create user",
		})
		return
	}
	tx.Commit()
}
