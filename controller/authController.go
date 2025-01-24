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

func Login(ctx *gin.Context) {
	var user model.UserLogin
	if err := ctx.ShouldBind(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "登陆失败",
		})
		return
	}

	var existedUser model.UserLogin
	if user.Phone != "" {
		if err := global.Db.Where("phone = ?", user.Phone).First(&existedUser).Error; err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"code":    http.StatusBadRequest,
				"message": "用户未注册",
			})
			return
		} else {
			if err := bcrypt.CompareHashAndPassword([]byte(existedUser.Password), []byte(user.Password)); err != nil {
				ctx.JSON(http.StatusBadRequest, gin.H{
					"code":    http.StatusBadRequest,
					"message": "手机号/邮箱或密码错误",
				})
				return
			}
			GetUserLoginInfo(existedUser.Uid, ctx)
		}
	}
	if user.Email != "" {
		if err := global.Db.Where("email = ?", user.Email).First(&existedUser).Error; err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"code":    http.StatusBadRequest,
				"message": "用户未注册",
			})
			return
		} else {
			if err := bcrypt.CompareHashAndPassword([]byte(existedUser.Password), []byte(user.Password)); err != nil {
				ctx.JSON(http.StatusBadRequest, gin.H{
					"code":    http.StatusBadRequest,
					"message": "手机号/邮箱或密码错误",
				})
				return
			}
			GetUserLoginInfo(existedUser.Uid, ctx)
		}
	}
}

func GetUserLoginInfo(uid uint, ctx *gin.Context) {
	var userInfo *model.UserInfo
	var userCreationInfo *model.UserCreationInfo
	if temp, err := getUserInfo(uid); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "登陆失败",
		})
		return
	} else {
		userInfo = temp
	}
	if temp, err := getUserCreationInfo(uid); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "登陆失败",
		})
		return
	} else {
		userCreationInfo = temp
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "登陆成功",
		"data": gin.H{
			"userInfo":         userInfo,
			"userCreationInfo": userCreationInfo,
		},
	})
}
