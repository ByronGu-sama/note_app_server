package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"note_app_server/repository"
	"note_app_server/response"
	"strconv"
)

// NewNote 创建笔记
func NewNote(ctx *gin.Context) {

}

// DelNote 删除笔记
func DelNote(ctx *gin.Context) {

}

// EditNote 编辑笔记
func EditNote(ctx *gin.Context) {

}

// GetNote 获取笔记
func GetNote(ctx *gin.Context) {
	tempNid := ctx.Param("nid")
	if tempNid == "" {
		response.RespondWithStatusBadRequest(ctx, "缺少关键信息")
		return
	}

	var nid uint
	if ans, err := strconv.Atoi(tempNid); err != nil {
		response.RespondWithStatusBadRequest(ctx, "转换错误")
		return
	} else {
		nid = uint(ans)
	}

	note, err := repository.GetNoteWithNid(nid)
	if err != nil {
		response.RespondWithStatusBadRequest(ctx, "无相关信息")
		return
	}
	if note.Status == 0 {
		ctx.JSON(http.StatusOK, gin.H{
			"code":    http.StatusOK,
			"message": "此条笔记已被删除/封禁",
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "success",
		"data":    note,
	})
}

// LikeNote 点赞笔记
func LikeNote(ctx *gin.Context) {

}

// DislikeNote 取消点赞笔记
func DislikeNote(ctx *gin.Context) {

}

// CollectNote 收藏笔记
func CollectNote(ctx *gin.Context) {

}

// CancelCollectNote 取消收藏笔记
func CancelCollectNote(ctx *gin.Context) {

}
