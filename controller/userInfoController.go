package controller

import (
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"net/http"
	"note_app_server/config"
	"note_app_server/model"
	"note_app_server/repository"
	"note_app_server/response"
	"note_app_server/service"
)

// GetAvatarUrl 获取代理头像地址
func GetAvatarUrl(ctx *gin.Context) {
	fileName := ctx.Param("fileName")
	reader, err := service.GetOssObject(config.AC.Oss.BucketName, fileName)

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

// GetUserInfo 获取用户详情
func GetUserInfo(ctx *gin.Context) {
	tempUid, ok := ctx.Get("uid")
	uid := tempUid.(uint)
	if !ok {
		response.RespondWithStatusBadRequest(ctx, "获取用户信息失败")
		return
	}

	var userInfo *model.UserInfo
	var userCreationInfo *model.UserCreationInfo

	if temp, err := repository.GetUserInfo(uid); err != nil {
		response.RespondWithStatusBadRequest(ctx, "获取用户信息失败")
		return
	} else {
		userInfo = temp
		userInfo.AvatarUrl = "http://localhost:8081/avatar/" + userInfo.AvatarUrl
	}

	if temp, err := repository.GetUserCreationInfo(uid); err != nil {
		response.RespondWithStatusBadRequest(ctx, "获取用户信息失败")
		return
	} else {
		userCreationInfo = temp
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"data": gin.H{
			"uid":       userInfo.Uid,
			"username":  userInfo.Username,
			"avatarUrl": userInfo.AvatarUrl,
			"age":       userInfo.Age,
			"gender":    userInfo.Gender,
			"birth":     userInfo.Birth,
			"signature": userInfo.Signature,
			"address":   userInfo.Address,
			"language":  userInfo.Language,
			"collects":  userCreationInfo.Collects,
			"followers": userCreationInfo.Followers,
			"follows":   userCreationInfo.Follows,
			"likes":     userCreationInfo.Likes,
			"noteCount": userCreationInfo.NoteCount,
		},
	})
}

// UpdateUserInfo 修改用户信息
func UpdateUserInfo(ctx *gin.Context) {
	tempUid, ok := ctx.Get("uid")
	uid := tempUid.(uint)
	if !ok {
		response.RespondWithStatusBadRequest(ctx, "获取用户信息失败")
		return
	}

	var userInfo *model.UserInfo
	if err := ctx.ShouldBind(&userInfo); err != nil {
		response.RespondWithStatusBadRequest(ctx, "绑定失败")
		return
	}
	userInfo.Uid = uid
	if err := repository.UpdateUserInfo(userInfo); err != nil {
		response.RespondWithStatusBadRequest(ctx, "更新失败")
		return
	}
	response.RespondWithStatusOK(ctx, "更新成功")
}

// UploadAvatar 上传头像
func UploadAvatar(ctx *gin.Context) {

}
