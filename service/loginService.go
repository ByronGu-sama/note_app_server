package service

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"note_app_server/global"
	"note_app_server/model/appModel"
	"time"
)

// CheckAccountStatus 判断账户状态
func CheckAccountStatus(status int64) error {
	if status == 0 {
		return errors.New("account has been banned")
	}
	if status == 2 {
		return errors.New("account has been canceled")
	}
	return nil
}

// GenerateJWT 用户登录后生成JWT 30天过期
func GenerateJWT(uid int64) (string, error) {
	mapClaims := jwt.MapClaims{
		"Iss": "note_app",
		"Sub": "token",
		"Uid": uid,
		"Exp": time.Now().Add(time.Hour * 24 * 30).Unix(),
		"Iat": time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, mapClaims)
	signedToken, err := token.SignedString(global.JWTKey)
	if err != nil {
		return "", err
	}
	return signedToken, nil
}

// ParseJWT 解析JWT
func ParseJWT(tokenString string) (interface{}, error) {
	claims := &appModel.JWT{}
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
