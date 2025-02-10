package service

import (
	"crypto/rand"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"log"
	"note_app_server/global"
	"note_app_server/model/appModel"
	"time"
)

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
func GenerateJWT(uid uint) (string, error) {
	mapClaims := jwt.MapClaims{
		"Iss": "note_app",
		"Sub": "token",
		"Uid": uid,
		"Exp": time.Now().Add(time.Hour * 24 * 30).Unix(),
		"Iat": time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, mapClaims)
	return token.SignedString(global.JWTKey)
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

// CreateJWTKey 生成JWT密钥
func CreateJWTKey() {
	var jwtKey = make([]byte, 32)
	if _, err := rand.Read(jwtKey); err != nil {
		log.Fatalf("Failed to generate JWT Key: %v", err)
	} else {
		global.JWTKey = jwtKey
	}
}
