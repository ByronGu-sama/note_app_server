package controller

import (
	"context"
	captcha20230305 "github.com/alibabacloud-go/captcha-20230305/client"
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

	if user.Password == "" {
		response.RespondWithStatusBadRequest(ctx, "密码不能为空")
		return
	}

	if len(user.Password) < 8 || len(user.Password) > 64 {
		response.RespondWithStatusBadRequest(ctx, "密码长度要求8～64个字符")
		return
	}

	// 加密密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		response.RespondWithStatusInternalServerError(ctx, "服务器内部错误")
		return
	}

	user.Password = string(hashedPassword)
	if err = repository.RegisterUser(&user); err != nil {
		response.RespondWithStatusBadRequest(ctx, err.Error())
		return
	}

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
		response.RespondWithStatusBadRequest(ctx, "该用户未登录")
		return
	}

	response.RespondWithStatusOK(ctx, "已退出登录")
}

// CheckToken token校验
func CheckToken(ctx *gin.Context) {
	response.RespondWithStatusOK(ctx, "valid")
}

// CheckCaptcha 验证码校验
func CheckCaptcha(ctx *gin.Context) {
	c := global.CaptchaClientPool.Get().(*captcha20230305.Client)
	defer global.CaptchaClientPool.Put(c)

	request := captcha20230305.VerifyIntelligentCaptchaRequest{}
	err := ctx.ShouldBind(&request)
	if err != nil {
		response.RespondWithStatusBadRequest(ctx, err.Error())
		return
	}

	resp, err1 := c.VerifyIntelligentCaptcha(&request)
	if err1 != nil {
		response.RespondWithStatusBadRequest(ctx, "校验失败")
		return
	}

	captchaVerifyResult := resp.Body.Result.VerifyResult
	captchaVerifyCode := resp.Body.Result.VerifyCode

	if *captchaVerifyResult {
		response.RespondWithStatusOK(ctx, "success")
	} else {
		response.RespondWithStatusBadRequest(ctx, *captchaVerifyCode)
	}
}
