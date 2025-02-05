package utils

import (
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
)

// EncodeNoteId 加密笔记id
func EncodeNoteId(name string) string {
	hash := md5.Sum([]byte(name))
	return hex.EncodeToString(hash[:])
}

// EncodePicsName 加密图片名称
func EncodePicsName(name string) string {
	hash := sha256.Sum256([]byte(name))
	return hex.EncodeToString(hash[:])
}
