package controller

import (
	"context"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"note_app_server/global"
	"note_app_server/model/userModel"
	"note_app_server/repository"
	"note_app_server/response"
	"note_app_server/service"
	"strconv"
)

// Register 注册
func Register(ctx *gin.Context) {
	var user userModel.UserLogin
	// 检查必要字段是否缺失
	if err := ctx.ShouldBind(&user); err != nil {
		response.RespondWithStatusBadRequest(ctx, "注册失败")
		return
	}

	//检查用户名或其他唯一字段是否已存在
	if _, err := repository.GetUserLoginInfoByPhone(user.Phone); err == nil {
		response.RespondWithStatusBadRequest(ctx, "该手机号已注册")
		return
	}

	// 检查密码是否为空
	if user.Password == "" {
		response.RespondWithStatusBadRequest(ctx, "密码不能为空")
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		response.RespondWithStatusInternalServerError(ctx, "服务器内部错误")
		return
	}

	user.Password = string(hashedPassword)
	// 使用事务创建用户登录信息
	tx := global.Db.Begin()
	if err := tx.Create(&user).Error; err != nil {
		tx.Rollback()
		response.RespondWithStatusInternalServerError(ctx, "服务器内部错误")
		return
	}
	tx.Commit()

	tx = global.Db.Begin()
	newUser, err := repository.GetUserLoginInfoByPhone(user.Phone) //获取系统生成的用户uid
	if err != nil {
		tx.Rollback()
		response.RespondWithStatusInternalServerError(ctx, "服务器内部错误")
		return
	}
	var userInfo userModel.UserInfo
	var userCreationInfo userModel.UserCreationInfo
	userInfo.Uid = newUser.Uid
	userInfo.AvatarUrl = "test.jpeg"
	userCreationInfo.Uid = newUser.Uid
	// 创建用户详细信息
	if err := tx.Create(&userInfo).Error; err != nil {
		tx.Rollback()
		response.RespondWithStatusInternalServerError(ctx, "服务器内部错误")
		return
	}
	// 创建用户创造者信息
	if err := tx.Create(&userCreationInfo).Error; err != nil {
		tx.Rollback()
		response.RespondWithStatusInternalServerError(ctx, "服务器内部错误")
		return
	}
	tx.Commit()
	response.RespondWithStatusOK(ctx, "注册成功")
}

// Login 登陆
func Login(ctx *gin.Context) {
	var user userModel.UserLogin
	if err := ctx.ShouldBind(&user); err != nil {
		response.RespondWithStatusBadRequest(ctx, "登陆失败")
		return
	}

	if user.Phone != "" {
		if existedUser, err := repository.GetUserLoginInfoByPhone(user.Phone); err != nil {
			response.RespondWithStatusBadRequest(ctx, "用户未注册")
			return
		} else {
			if err1 := service.CheckAccountStatus(existedUser.AccountStatus); err1 != nil {
				response.RespondWithStatusBadRequest(ctx, "用户状态异常")
				return
			}
			if err := bcrypt.CompareHashAndPassword([]byte(existedUser.Password), []byte(user.Password)); err != nil {
				repository.UpdateLoginFailedAt(existedUser.Uid)
				response.RespondWithStatusBadRequest(ctx, "手机号/密码错误")
				return
			}
			if token, err := repository.GetToken(existedUser.Uid); err != nil {
				response.RespondWithStatusBadRequest(ctx, "登陆失败")
			} else {
				ctx.JSON(http.StatusOK, gin.H{
					"code":    http.StatusOK,
					"message": "登陆成功",
					"token":   token,
				})
			}
		}
	} else {
		response.RespondWithStatusBadRequest(ctx, "手机号不能为空")
		return
	}
}

// Logout 登出
func Logout(ctx *gin.Context) {
	tempUid, ok := ctx.Get("uid")
	uid := tempUid.(uint)
	if !ok {
		response.RespondWithStatusBadRequest(ctx, "获取用户信息失败")
		return
	}

	rCtx := context.Background()
	deletedCount, err := global.TokenRdb.Del(rCtx, strconv.Itoa(int(uid))).Result()
	if err != nil {
		response.RespondWithStatusBadRequest(ctx, "退出登录失败")
		return
	}
	if deletedCount == 0 {
		response.RespondWithStatusBadRequest(ctx, "该用户未登录或已退出")
		return
	}

	response.RespondWithStatusOK(ctx, "已退出登录")
}

// CheckToken 检查token有效性
func CheckToken(ctx *gin.Context) {
	response.RespondWithStatusOK(ctx, "valid")
}
