package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"math/rand/v2"
	"net/http"
	"note_app_server/model/commentModel"
	"note_app_server/repository"
	"note_app_server/response"
	"note_app_server/utils"
	"strconv"
	"time"
)

// NewComment 发送评论
func NewComment(ctx *gin.Context) {
	uid, ok := ctx.Get("uid")
	if !ok {
		response.RespondWithStatusBadRequest(ctx, "获取uid失败")
		return
	}
	var cmt *commentModel.Comment
	if err := ctx.ShouldBind(&cmt); err != nil {
		response.RespondWithStatusBadRequest(ctx, "绑定失败")
		return
	}

	if cmt.Nid == "" || cmt.Content == "" {
		response.RespondWithStatusBadRequest(ctx, "缺少信息")
		return
	}

	cid := utils.EncodeWithSHA256(fmt.Sprintf("%s-%d-%d", cmt.Nid, uid, rand.Int64()))
	cmt.Cid = cid
	cmt.Uid = uid.(uint)
	cmt.CreatedAt = time.Now()
	if cmt.RootId == "" {
		cmt.RootId = cid
	}

	cmtInfo := &commentModel.CommentsInfo{Cid: cid, LikesCount: 0}

	result, err := repository.NewComment(cmt, cmtInfo)
	if err != nil {
		response.RespondWithStatusBadRequest(ctx, "评论失败")
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "获取成功",
		"data":    result,
	})
}

// DelComment 删除评论
func DelComment(ctx *gin.Context) {
	uid, ok := ctx.Get("uid")
	if !ok {
		response.RespondWithStatusBadRequest(ctx, "获取uid失败")
		return
	}
	cid := ctx.Param("cid")
	if cid == "" {
		response.RespondWithStatusBadRequest(ctx, "缺少信息")
		return
	}

	if err := repository.DeleteComment(uid.(uint), cid); err != nil {
		response.RespondWithStatusBadRequest(ctx, err.Error())
		return
	}

	response.RespondWithStatusOK(ctx, "success")
}

// LikeComment 点赞评论
func LikeComment(ctx *gin.Context) {
	uid, ok := ctx.Get("uid")
	if !ok {
		response.RespondWithStatusBadRequest(ctx, "获取uid失败")
		return
	}
	cid := ctx.Param("cid")
	if cid == "" {
		response.RespondWithStatusBadRequest(ctx, "缺少信息")
		return
	}

	if err := repository.LikeComment(uid.(uint), cid); err != nil {
		response.RespondWithStatusBadRequest(ctx, err.Error())
		return
	}
	response.RespondWithStatusOK(ctx, "点赞成功")
}

// CancelLikeComment 取消点赞评论
func CancelLikeComment(ctx *gin.Context) {
	uid, ok := ctx.Get("uid")
	if !ok {
		response.RespondWithStatusBadRequest(ctx, "获取uid失败")
		return
	}
	cid := ctx.Param("cid")
	if cid == "" {
		response.RespondWithStatusBadRequest(ctx, "缺少信息")
		return
	}

	if err := repository.DislikeComment(uid.(uint), cid); err != nil {
		response.RespondWithStatusBadRequest(ctx, err.Error())
		return
	}
	response.RespondWithStatusOK(ctx, "取消点赞成功")
}

// GetCommentList 获取评论列表
func GetCommentList(ctx *gin.Context) {
	nid := ctx.Param("nid")
	page := ctx.Query("page")
	limit := ctx.Query("limit")
	if nid == "" || page == "" || limit == "" {
		response.RespondWithStatusBadRequest(ctx, "缺少信息")
		return
	}
	truePage, err1 := strconv.Atoi(page)
	if err1 != nil {
		response.RespondWithStatusBadRequest(ctx, "参数错误")
		return
	}
	trueLimit, err2 := strconv.Atoi(limit)
	if err2 != nil {
		response.RespondWithStatusBadRequest(ctx, "参数错误")
		return
	}
	commentList, err := repository.GetNoteCommentsList(nid, truePage, trueLimit)
	if err != nil {
		return
	}

	rootArr := make([]commentModel.CommentDetail, 0)
	for _, v := range commentList {
		v.AvatarUrl = utils.AddAvatarPrefix(v.AvatarUrl)
		rootArr = append(rootArr, v)
	}
	ctx.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "success",
		"data":    rootArr,
	})
}

// GetSubCommentList 获取子评论列表
func GetSubCommentList(ctx *gin.Context) {
	nid := ctx.Param("nid")
	rootId := ctx.Param("rootId")
	page := ctx.Query("page")
	limit := ctx.Query("limit")
	if nid == "" || rootId == "" || page == "" || limit == "" {
		response.RespondWithStatusBadRequest(ctx, "缺少信息")
		return
	}
	truePage, err1 := strconv.Atoi(page)
	if err1 != nil {
		response.RespondWithStatusBadRequest(ctx, "参数错误")
		return
	}
	trueLimit, err2 := strconv.Atoi(limit)
	if err2 != nil {
		response.RespondWithStatusBadRequest(ctx, "参数错误")
		return
	}
	commentList, err := repository.GetSubCommentsList(nid, rootId, truePage, trueLimit)
	if err != nil {
		return
	}

	rootArr := make([]commentModel.CommentDetail, 0)
	for _, v := range commentList {
		v.AvatarUrl = utils.AddAvatarPrefix(v.AvatarUrl)
		rootArr = append(rootArr, v)
	}
	ctx.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "success",
		"data":    rootArr,
	})
}
