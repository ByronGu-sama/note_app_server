package controller

import (
	"context"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"note_app_server1/global"
	"note_app_server1/model"
	"note_app_server1/repository"
	"note_app_server1/service"
	"strconv"
)

// Register 注册
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
	if _, err := repository.GetUserLoginInfoByPhone(user.Phone); err == nil {
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

	// 使用事务创建用户登录信息
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
	newUser, err := repository.GetUserLoginInfoByPhone(user.Phone) //获取系统生成的用户uid
	if err != nil {
		tx.Rollback()
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to create user",
		})
		return
	}
	var userInfo model.UserInfo
	var userCreationInfo model.UserCreationInfo
	userInfo.Uid = newUser.Uid
	userCreationInfo.Uid = newUser.Uid
	// 创建用户详细信息
	if err := tx.Create(&userInfo).Error; err != nil {
		tx.Rollback()
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to create user",
		})
		return
	}
	// 创建用户创造者信息
	if err := tx.Create(&userCreationInfo).Error; err != nil {
		tx.Rollback()
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to create user creation information",
		})
		return
	}
	tx.Commit()

	ctx.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "注册成功",
	})
}

// Login 登陆
func Login(ctx *gin.Context) {
	var user model.UserLogin
	if err := ctx.ShouldBind(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "登陆失败",
		})
		return
	}

	if user.Phone != "" {
		if existedUser, err := repository.GetUserLoginInfoByPhone(user.Phone); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"code":    http.StatusBadRequest,
				"message": "用户未注册",
			})
			return
		} else {
			if err := service.CheckAccountStatus(existedUser.AccountStatus, ctx); err != nil {
				return
			}
			if err := bcrypt.CompareHashAndPassword([]byte(existedUser.Password), []byte(user.Password)); err != nil {
				repository.UpdateLoginFailedAt(existedUser.Uid)
				ctx.JSON(http.StatusBadRequest, gin.H{
					"code":    http.StatusBadRequest,
					"message": "手机号/邮箱或密码错误",
				})
				return
			}
			service.GetUserLoginInfo(existedUser.Uid, ctx)
		}
	}
	if user.Email != "" {
		if existedUser, err := repository.GetUserLoginInfoByEmail(user.Email); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"code":    http.StatusBadRequest,
				"message": "用户未注册",
			})
			return
		} else {
			if err := service.CheckAccountStatus(existedUser.AccountStatus, ctx); err != nil {
				return
			}
			if err := bcrypt.CompareHashAndPassword([]byte(existedUser.Password), []byte(user.Password)); err != nil {
				repository.UpdateLoginFailedAt(existedUser.Uid)
				ctx.JSON(http.StatusBadRequest, gin.H{
					"code":    http.StatusBadRequest,
					"message": "手机号/邮箱或密码错误",
				})
				return
			}
			service.GetUserLoginInfo(existedUser.Uid, ctx)
		}
	}
}

// Logout 登出
func Logout(ctx *gin.Context) {
	var user model.UserInfo
	token := ctx.GetHeader("token")

	if token == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "无权限",
		})
		return
	}

	if err := ctx.ShouldBind(&user); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "校验失败",
		})
		return
	}

	if _, err := service.ParseJWT(token); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "无权限",
		})
	}

	rCtx := context.Background()
	if err := global.TokenRdb.Del(rCtx, strconv.Itoa(int(user.Uid))).Err(); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "退出登录失败",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "已退出登录",
	})
}
