package controller

import (
	"errors"
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

	if note.Title == "" || note.Content == "" {
		response.RespondWithStatusBadRequest(ctx, "关键信息不能为空")
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
	if err := tx.Model(userCreation).Where("uid = ?", uid).Update("noteCount", gorm.Expr("noteCount + ?", 1)).Error; err != nil {
		tx.Rollback()
		response.RespondWithStatusBadRequest(ctx, "创建失败")
		return
	}
	tx.Commit()
	response.RespondWithStatusOK(ctx, "创建成功")
}

// DelNote 删除笔记
func DelNote(ctx *gin.Context) {
	nid, uid, err := checkUidAndNid(ctx)
	if err != nil {
		response.RespondWithStatusBadRequest(ctx, err.Error())
		return
	}
	if err := repository.DeleteNoteWithUid(nid, uid); err != nil {
		response.RespondWithStatusBadRequest(ctx, "删除失败")
		return
	}
	response.RespondWithStatusOK(ctx, "删除成功")
}

// EditNote 编辑笔记
func EditNote(ctx *gin.Context) {
	var note model.Note
	if err := ctx.ShouldBind(&note); err != nil {
		response.RespondWithStatusBadRequest(ctx, "绑定失败")
		return
	}

	tempUid, ok := ctx.Get("uid")
	uid := tempUid.(uint)
	if !ok {
		response.RespondWithStatusBadRequest(ctx, "获取用户信息失败")
		return
	}
	note.Uid = uid
	if err := repository.UpdateNoteWithUid(&note); err != nil {
		response.RespondWithStatusBadRequest(ctx, "更新失败")
		return
	}
	response.RespondWithStatusOK(ctx, "更新成功")
}

// GetNote 获取笔记
func GetNote(ctx *gin.Context) {
	nid := ctx.Param("nid")
	if nid == "" {
		response.RespondWithStatusBadRequest(ctx, "缺少关键信息")
		return
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
	nid, uid, err := checkUidAndNid(ctx)
	if err != nil {
		response.RespondWithStatusBadRequest(ctx, err.Error())
		return
	}

	if err := repository.LikeNote(nid, uid); err != nil {
		response.RespondWithStatusBadRequest(ctx, "点赞失败")
		return
	}
	response.RespondWithStatusOK(ctx, "点赞成功")
}

// DislikeNote 取消点赞笔记
func DislikeNote(ctx *gin.Context) {
	nid, uid, err := checkUidAndNid(ctx)
	if err != nil {
		response.RespondWithStatusBadRequest(ctx, err.Error())
		return
	}

	if err := repository.CancelLikeNote(nid, uid); err != nil {
		response.RespondWithStatusBadRequest(ctx, "取消点赞失败")
		return
	}
	response.RespondWithStatusOK(ctx, "取消点赞成功")
}

// CollectNote 收藏笔记
func CollectNote(ctx *gin.Context) {
	nid, uid, err := checkUidAndNid(ctx)
	if err != nil {
		response.RespondWithStatusBadRequest(ctx, err.Error())
		return
	}

	if err := repository.CollectNote(nid, uid); err != nil {
		response.RespondWithStatusBadRequest(ctx, "收藏失败")
		return
	}
	response.RespondWithStatusOK(ctx, "收藏成功")
}

// CancelCollectNote 取消收藏笔记
func CancelCollectNote(ctx *gin.Context) {
	nid, uid, err := checkUidAndNid(ctx)
	if err != nil {
		response.RespondWithStatusBadRequest(ctx, err.Error())
		return
	}

	if err := repository.CancelCollectNote(nid, uid); err != nil {
		response.RespondWithStatusBadRequest(ctx, "取消收藏失败")
		return
	}
	response.RespondWithStatusOK(ctx, "取消收藏成功")
}

// GetNoteList 获取笔记列表
func GetNoteList(ctx *gin.Context) {
	tempPage := ctx.Query("page")
	tempLimit := ctx.Query("limit")

	if tempPage == "" || tempLimit == "" {
		response.RespondWithStatusBadRequest(ctx, "缺失参数")
		return
	}
	page, err1 := strconv.Atoi(tempPage)
	if err1 != nil {
		response.RespondWithStatusBadRequest(ctx, err1.Error())
		return
	}
	limit, err2 := strconv.Atoi(tempLimit)
	if err2 != nil {
		response.RespondWithStatusBadRequest(ctx, err2.Error())
		return
	}

	result, err3 := repository.GetNoteList(page, limit)
	if err3 != nil {
		response.RespondWithStatusBadRequest(ctx, err3.Error())
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "success",
		"data":    result,
	})
}

func checkUidAndNid(ctx *gin.Context) (string, uint, error) {
	nid := ctx.Param("nid")
	if nid == "" {
		return "", 0, errors.New("缺少必要信息")
	}

	tempUid, ok := ctx.Get("uid")
	uid := tempUid.(uint)
	if !ok {
		return "", 0, errors.New("获取uid失败")
	}
	return nid, uid, nil
}
