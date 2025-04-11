package response

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func RespondWithUnauthorized(ctx *gin.Context, msg string) {
	ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
		"code":    http.StatusUnauthorized,
		"message": msg,
	})
}

func RespondWithStatusOK(ctx *gin.Context, msg string) {
	ctx.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": msg,
	})
}

func RespondWithStatusBadRequest(ctx *gin.Context, msg string) {
	ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
		"code":    http.StatusBadRequest,
		"message": msg,
	})
}

func RespondWithStatusInternalServerError(ctx *gin.Context, msg string) {
	ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
		"code":    http.StatusInternalServerError,
		"message": msg,
	})
}

func RespondWithStatusServiceUnavailable(ctx *gin.Context, msg string) {
	ctx.AbortWithStatusJSON(http.StatusServiceUnavailable, gin.H{
		"code":    http.StatusServiceUnavailable,
		"message": msg,
	})
}
