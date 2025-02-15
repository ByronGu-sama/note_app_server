package controller

import (
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"net/http"
	"note_app_server/config"
	"note_app_server/repository"
	"note_app_server/response"
	"note_app_server/service"
	"note_app_server/utils"
)

// GetProfileBannerUrl 获取banner代理地址
func GetProfileBannerUrl(ctx *gin.Context) {
	fileName := ctx.Param("bid")
	reader, err := service.GetOssObject(config.AC.Oss.StyleBucket, "profileBanner/", fileName)

	if err != nil {
		response.RespondWithStatusBadRequest(ctx, "获取Oss服务失败")
		return
	}

	defer func(reader io.ReadCloser) {
		err := reader.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(reader)

	data, err := io.ReadAll(reader)
	if err != nil {
		response.RespondWithStatusBadRequest(ctx, "读取文件流失败")
		return
	}

	ctx.Header("Content-Type", "image/jpeg")
	ctx.Data(http.StatusOK, "image/jpeg", data)
}

// UpdateProfileBanner 更新banner
func UpdateProfileBanner(ctx *gin.Context) {
	tempUid, ok := ctx.Get("uid")
	uid := tempUid.(uint)
	if !ok {
		response.RespondWithStatusBadRequest(ctx, "获取用户信息失败")
		return
	}

	bannerFile, err1 := ctx.FormFile("file")
	if err1 != nil {
		response.RespondWithStatusBadRequest(ctx, err1.Error())
		return
	}

	openBanner, err2 := bannerFile.Open()
	defer openBanner.Close()

	if err2 != nil {
		response.RespondWithStatusBadRequest(ctx, err2.Error())
		return
	}

	// 读取文件后重置指针
	tempFile, err := io.ReadAll(openBanner)
	if err != nil {
		return
	}
	_, err3 := openBanner.Seek(0, io.SeekStart)
	if err3 != nil {
		response.RespondWithStatusBadRequest(ctx, err3.Error())
		return
	}

	fileType := utils.DetectFileType(tempFile)

	bannerName, err4 := service.UploadFileObject(config.AC.Oss.StyleBucket, "profileBanner/", openBanner, fileType)
	if err4 != nil {
		response.RespondWithStatusBadRequest(ctx, err4.Error())
		return
	}

	lastBanner, err5 := repository.GetLastBanner(uid)
	if err5 != nil {
		log.Fatal(err5)
	}

	if err6 := repository.UpdateProfileBanner(uid, bannerName); err6 != nil {
		response.RespondWithStatusInternalServerError(ctx, err6.Error())
		err7 := service.DeleteObject(config.AC.Oss.StyleBucket, "profileBanner/", bannerName)
		if err7 != nil {
			log.Fatal(err7)
		}
		return
	}

	err8 := service.DeleteObject(config.AC.Oss.StyleBucket, "profileBanner/", lastBanner)
	if err8 != nil {
		log.Fatal(err8)
	}

	response.RespondWithStatusOK(ctx, "上传成功")
}

// GetStyle 获取app风格
func GetStyle(ctx *gin.Context) {
	tempUid, ok := ctx.Get("uid")
	uid := tempUid.(uint)
	if !ok {
		response.RespondWithStatusBadRequest(ctx, "获取用户信息失败")
		return
	}
	result, err := repository.GetStyle(uid)
	if err != nil {
		response.RespondWithStatusBadRequest(ctx, err.Error())
		return
	}

	result.ProfileBanner = utils.AddProfileBannerPrefix(result.ProfileBanner)
	ctx.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "获取信息成功",
		"data":    result,
	})
}
