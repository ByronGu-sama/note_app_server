package controller

import (
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"net/http"
	"note_app_server1/config"
	"note_app_server1/service"
)

// GetAvatarUrl 获取代理头像地址
func GetAvatarUrl(ctx *gin.Context) {
	fileName := ctx.Param("fileName")
	reader, err := service.GetOssObject(config.AC.Oss.BucketName, fileName)

	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"code":    http.StatusBadRequest,
			"message": err.Error(),
		})
	}

	defer func(reader io.ReadCloser) {
		err := reader.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(reader)

	data, err := io.ReadAll(reader)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"code":    http.StatusBadRequest,
			"message": err.Error(),
		})
	}

	ctx.Header("Content-Type", "image/jpeg")
	ctx.Data(http.StatusOK, "image/jpeg", data)
}
