package controller

import (
	"github.com/gin-gonic/gin"
	"note_app_server/service"
)

func GetWS(ctx *gin.Context) {
	service.InitWS(ctx)
}
