package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"math/rand/v2"
	"note_app_server/model/commentModel"
	"note_app_server/repository"
	"note_app_server/response"
	"note_app_server/utils"
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

	cmtInfo := &commentModel.CommentsInfo{cid, 0}

	err := repository.NewComment(cmt, cmtInfo)
	if err != nil {
		response.RespondWithStatusBadRequest(ctx, "评论失败")
		return
	}
	response.RespondWithStatusOK(ctx, "已发送评论")
}

// DelComment 删除评论
func DelComment(ctx *gin.Context) {
	uid, ok := ctx.Get("uid")
	if !ok {
		response.RespondWithStatusBadRequest(ctx, "获取uid失败")
		return
	}
	cid := ctx.Param("cid")
	nid := ctx.Param("nid")
	if cid == "" || nid == "" {
		response.RespondWithStatusBadRequest(ctx, "缺少信息")
		return
	}

	if err := repository.DeleteComment(uid.(uint), cid, nid); err != nil {
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

}
