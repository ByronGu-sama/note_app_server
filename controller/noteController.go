package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"math/rand"
	"net/http"
	"note_app_server/global"
	"note_app_server/model"
	"note_app_server/repository"
	"note_app_server/response"
	"note_app_server/utils"
	"strconv"
	"time"
)

// NewNote 创建笔记
func NewNote(ctx *gin.Context) {
	var note model.Note
	if err := ctx.ShouldBind(&note); err != nil {
		response.RespondWithStatusBadRequest(ctx, "创建失败")
		return
	}

	tempUid, ok := ctx.Get("uid")
	uid := tempUid.(uint)
	if !ok {
		response.RespondWithStatusBadRequest(ctx, "获取用户信息失败")
		return
	}
	note.Uid = uid
	noteName := utils.EncodeName(fmt.Sprintf("%d-%d-%d", time.Now().Unix(), uid, rand.Int63()))
	note.Nid = noteName

	tx := global.Db.Begin()
	if err := tx.Create(&note).Error; err != nil {
		tx.Rollback()
		response.RespondWithStatusBadRequest(ctx, "创建失败")
		return
	}
	userCreation := &model.UserCreationInfo{}
	if err := tx.Model(userCreation).Where("uid = ?", uid).Update("noteCount", gorm.Expr("noteCount+ ?", 1)).Error; err != nil {
		tx.Rollback()
		response.RespondWithStatusBadRequest(ctx, "创建失败")
		return
	}
	tx.Commit()
	response.RespondWithStatusOK(ctx, "创建成功")
}

// DelNote 删除笔记
func DelNote(ctx *gin.Context) {
	tempUid, ok := ctx.Get("uid")
	uid := tempUid.(uint)
	if !ok {
		response.RespondWithStatusBadRequest(ctx, "获取用户信息失败")
		return
	}

	nid := ctx.Param("nid")
	if err := repository.DeleteNoteWithUid(nid, uid); err != nil {
		response.RespondWithStatusBadRequest(ctx, "删除失败")
		return
	}
	response.RespondWithStatusOK(ctx, "删除成功")
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
