package utils

import (
	"crypto/md5"
	"encoding/hex"
)

// EncodeName 加密名称
func EncodeName(name string) string {
	hash := md5.Sum([]byte(name))
	return hex.EncodeToString(hash[:])
}
