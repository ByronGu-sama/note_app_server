package appModel

import (
	"github.com/golang-jwt/jwt/v5"
)

type JWT struct {
	Iss string // 发行者
	Sub string // 主题
	Uid uint   // uid
	Exp int64  // 过期时间
	Iat int64  // 发行时间
	jwt.RegisteredClaims
}
