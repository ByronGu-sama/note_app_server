package utils

import (
	"bytes"
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"github.com/gabriel-vasile/mimetype"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"image/png"
	"io"
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

// DetectFileType 返回文件类型
func DetectFileType(file []byte) (string, error) {
	fileType := mimetype.Detect(file).String()
	var err error

	switch fileType {
	case "image/png":
		fileType = "png"
	case "image/jpeg":
		fileType = "jpeg"
	case "image/webp":
		fileType = "webp"
	case "image/heic":
		fileType = "heic"
	case "image/heif":
		fileType = "gif"
	default:
		err = errors.New("not pre defined file type")
	}
	return fileType, err
}

// AddAvatarPrefix 添加前端访问头像url前缀
func AddAvatarPrefix(url string) string {
	return "http://" + config.AC.App.Host + config.AC.App.Port + "/avatar/" + url
}

// AddNotePicPrefix 添加前端访问笔记图片url前缀
func AddNotePicPrefix(nid, url string) string {
	return "http://" + config.AC.App.Host + config.AC.App.Port + "/note/pic/" + nid + "/" + url
}

// AddProfileBannerPrefix 添加前端访问用户页banner前缀
func AddProfileBannerPrefix(url string) string {
	return "http://" + config.AC.App.Host + config.AC.App.Port + "/style/profileBanner/" + url
}

// CompressJPEGPic 压缩jpeg图片
func CompressJPEGPic(file io.Reader, quality int) ([]byte, error) {
	img, err := jpeg.Decode(file)
	if err != nil {
		return nil, err
	}
	buf := new(bytes.Buffer)
	err = jpeg.Encode(buf, img, &jpeg.Options{Quality: quality})
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// CompressPNGPic 压缩png图片
func CompressPNGPic(file io.Reader, quality int) ([]byte, error) {
	img, err := png.Decode(file)
	if err != nil {
		return nil, err
	}
	newImg := image.NewRGBA(img.Bounds())
	draw.Draw(newImg, newImg.Bounds(), &image.Uniform{C: color.White}, image.Point{}, draw.Src)
	draw.Draw(newImg, newImg.Bounds(), img, img.Bounds().Min, draw.Over)

	buf := new(bytes.Buffer)
	err = jpeg.Encode(buf, newImg, &jpeg.Options{Quality: quality})
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
