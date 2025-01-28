package utils

import (
	"crypto/sha256"
	"encoding/hex"
)

// EncodeFileName 加密文件名
func EncodeFileName(fileName string) string {
	hash := sha256.Sum256([]byte(fileName))
	return hex.EncodeToString(hash[:])
}
