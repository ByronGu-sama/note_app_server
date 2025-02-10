package utils

import (
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"mime/multipart"
	"net/http"
	"note_app_server/config"
)

// EncodeWithMD5 加密笔记id
func EncodeWithMD5(name string) string {
	hash := md5.Sum([]byte(name))
	return hex.EncodeToString(hash[:])
}

// EncodeWithSHA256 加密图片名称
func EncodeWithSHA256(name string) string {
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

// AddAvatarPrefix 添加头像url前缀
func AddAvatarPrefix(url string) string {
	return "http://" + config.AC.App.Host + config.AC.App.Port + "/avatar/" + url
}

// AddNotePicPrefix 添加笔记图片url前缀
func AddNotePicPrefix(nid, url string) string {
	return "http://" + config.AC.App.Host + config.AC.App.Port + "/note/pic/" + nid + "/" + url
}
