package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"note_app_server/config"
	"note_app_server/response"
	"note_app_server/service"
	"note_app_server/utils"
)

func UploadUserAvatar(ctx *gin.Context) {
	// 获取表单文件
	file, err := ctx.FormFile("avatar")
	if err != nil {
		response.RespondWithStatusBadRequest(ctx, err.Error())
		return
	}
	// 转换为file类型的文件
	openFile, err := file.Open()
	if err != nil {
		response.RespondWithStatusBadRequest(ctx, err.Error())
		return
	}

	// 检查文件类型
	contentType, err := utils.DetectFileType(&openFile)
	if err != nil {
		response.RespondWithStatusBadRequest(ctx, err.Error())
		return
	}
	if contentType != "image/jpeg" && contentType != "image/png" {
		response.RespondWithStatusBadRequest(ctx, "不支持的文件类型")
		return
	}
	if contentType == "image/png" {
		contentType = "png"
	}
	if contentType == "image/jpeg" {
		contentType = "jpeg"
	}

	// 上传文件至oss
	fileName, err1 := service.UploadFileObject(config.AC.Oss.BucketName, "tempAvatar/", openFile, contentType)
	if err1 != nil {
		response.RespondWithStatusBadRequest(ctx, err1.Error())
		return
	}
	// 响应
	ctx.JSON(http.StatusOK, gin.H{
		"code":      http.StatusOK,
		"message":   "上传成功",
		"avatarUrl": utils.AddAvatarPrefix(fileName),
	})
}
