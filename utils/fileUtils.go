package utils

import (
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"mime/multipart"
	"net/http"
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

// DetectFileType 检查文件类型
func DetectFileType(file *multipart.File) (string, error) {
	// 读取文件前512字节的数据
	buf := make([]byte, 512)
	if _, err := (*file).Read(buf); err != nil {
		return "", err
	}
	// 重置文件读取指针
	if _, err := (*file).Seek(0, io.SeekStart); err != nil {
		return "", err
	}
	// 检查文件类型
	contentType := http.DetectContentType(buf)
	return contentType, nil
}
