package service

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"note_app_server1/global"
	"note_app_server1/model"
	"note_app_server1/repository"
	"strconv"
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

	// 生成jwt并保存
	token, err := GenerateJWT(userInfo)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "登陆失败",
		})
	}

	rCtx := context.Background()
	err = global.TokenRdb.Set(rCtx, strconv.Itoa(int(userInfo.Uid)), token, time.Hour*24*30).Err()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "登陆失败",
		})
	}
	repository.UpdateLoginSuccessAt(userInfo.Uid)
	ctx.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "登陆成功",
		"token":   token,
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

// GenerateJWT 用户登录后生成JWT 30天过期
func GenerateJWT(user *model.UserInfo) (string, error) {
	mapClaims := jwt.MapClaims{
		"iss": "note_app",
		"sub": "token",
		"uid": user.Uid,
		"exp": time.Now().Add(time.Hour * 24 * 30).Unix(),
		"iat": time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, mapClaims)
	return token.SignedString(global.JWTKey)
}

// ParseJWT 解析JWT
func ParseJWT(tokenString string) (interface{}, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return global.JWTKey, nil
	})
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}
	return token.Claims, nil
}
