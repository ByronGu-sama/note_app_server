package service

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"log"
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
		ctx.JSON(http.StatusOK, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "登陆失败",
		})
		return
	} else {
		userInfo = temp
		userInfo.AvatarUrl = "http://localhost:8081/avatar/" + userInfo.AvatarUrl
	}
	if temp, err := repository.GetUserCreationInfo(uid); err != nil {
		ctx.JSON(http.StatusOK, gin.H{
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
		ctx.JSON(http.StatusOK, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "登陆失败",
		})
	}

	rCtx := context.Background()
	err = global.TokenRdb.Set(rCtx, strconv.Itoa(int(userInfo.Uid)), token, time.Hour*24*30).Err()
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
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
			"uid":       userInfo.Uid,
			"username":  userInfo.Username,
			"avatarUrl": userInfo.AvatarUrl,
			"age":       userInfo.Age,
			"gender":    userInfo.Gender,
			"birth":     userInfo.Birth,
			"signature": userInfo.Signature,
			"address":   userInfo.Address,
			"language":  userInfo.Language,
			"collects":  userCreationInfo.Collects,
			"followers": userCreationInfo.Followers,
			"follows":   userCreationInfo.Follows,
			"likes":     userCreationInfo.Likes,
			"noteCount": userCreationInfo.NoteCount,
		},
	})
}

// CheckAccountStatus 判断账户状态
func CheckAccountStatus(status uint) error {
	if status == 0 {
		return errors.New("account has been banned")
	}
	if status == 2 {
		return errors.New("account has been canceled")
	}
	return nil
}

// GenerateJWT 用户登录后生成JWT 30天过期
func GenerateJWT(user *model.UserInfo) (string, error) {
	mapClaims := jwt.MapClaims{
		"Iss": "note_app",
		"Sub": "token",
		"Uid": user.Uid,
		"Exp": time.Now().Add(time.Hour * 24 * 30).Unix(),
		"Iat": time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, mapClaims)
	return token.SignedString(global.JWTKey)
}

// ParseJWT 解析JWT
func ParseJWT(tokenString string) (interface{}, error) {
	claims := &model.JWT{}
	temp, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return global.JWTKey, nil
	})
	if err != nil {
		return nil, err
	}
	if !temp.Valid {
		return nil, fmt.Errorf("invalid token")
	}
	return claims, nil
}

// CreateJWTKey 生成JWT密钥
func CreateJWTKey() {
	var jwtKey = make([]byte, 32)
	if _, err := rand.Read(jwtKey); err != nil {
		log.Fatalf("Failed to generate JWT Key: %v", err)
	} else {
		global.JWTKey = jwtKey
	}
}
