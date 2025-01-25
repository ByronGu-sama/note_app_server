package service

import (
	"crypto/rand"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"note_app_server1/model"
	"note_app_server1/repository"
	"time"
)

// GetUserLoginInfo 获取用户相关的信息
func GetUserLoginInfo(uid uint, ctx *gin.Context) {
	var userInfo *model.UserInfo
	var userCreationInfo *model.UserCreationInfo
	if temp, err := repository.GetUserInfo(uid); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "登陆失败",
		})
		return
	} else {
		userInfo = temp
	}
	if temp, err := repository.GetUserCreationInfo(uid); err != nil {
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

// CheckAccountStatus 判断账户状态
func CheckAccountStatus(status uint, ctx *gin.Context) error {
	if status == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "账号已封禁，暂时无法登陆",
		})
		return fmt.Errorf("account has been banned")
	}
	if status == 2 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "账号已注销",
		})
		return fmt.Errorf("account has been canceled")
	}
	return nil
}

// GenerateJWT 用户登录后生成JWT
func GenerateJWT(user *model.UserLogin) (string, error) {
	key := make([]byte, 32)
	if _, err := rand.Read(key); err != nil {
		return "", err
	}
	mapClaims := jwt.MapClaims{
		"iss": "note_app",
		"sub": "token",
		"aud": user.Uid,
		"exp": time.Now().Add(time.Hour * 24 * 30).Unix(),
		"iat": time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, mapClaims)
	return token.SignedString(key)
}
