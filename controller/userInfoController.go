package controller

import (
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"net/http"
	"note_app_server/config"
	"note_app_server/model/userModel"
	"note_app_server/repository"
	"note_app_server/response"
	"note_app_server/service"
	"note_app_server/utils"
)

// GetAvatarUrl 获取代理头像地址
func GetAvatarUrl(ctx *gin.Context) {
	fileName := ctx.Param("fileName")
	reader, err := service.GetOssObject(config.AC.Oss.BucketName, "avatar/", fileName)

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

// ChangeAvatar 确定修改头像
func ChangeAvatar(ctx *gin.Context) {
	// 获取url
	avatarUrl := ctx.Query("avatarUrl")
	if avatarUrl == "" {
		response.RespondWithStatusBadRequest(ctx, "头像地址为空")
		return
	}
	uid, _ := ctx.Get("uid")
	// 判断用户是否已上传头像
	exist, err1 := service.HasObject(config.AC.Oss.BucketName, "tempAvatar/", avatarUrl)
	if err1 != nil {
		response.RespondWithStatusBadRequest(ctx, err1.Error())
		return
	}
	if exist {
		// 将头像从temp文件夹转移至常规文件夹
		err := service.CopyObjectToAnother(config.AC.Oss.BucketName, "tempAvatar/"+avatarUrl, "avatar/"+avatarUrl)
		if err != nil {
			response.RespondWithStatusBadRequest(ctx, err.Error())
			return
		}
		// 从temp文件夹中删除文件
		err = service.DeleteObject(config.AC.Oss.BucketName, "tempAvatar/", avatarUrl)
		if err != nil {
			response.RespondWithStatusBadRequest(ctx, err.Error())
			return
		}
		// 更新用户头像地址
		err = repository.UpdateUserAvatar(uid.(uint), avatarUrl)
		if err != nil {
			response.RespondWithStatusInternalServerError(ctx, err.Error())
			return
		}

		response.RespondWithStatusOK(ctx, "保存成功")
	} else {
		response.RespondWithStatusBadRequest(ctx, "头像未上传")
		return
	}
}

// GetUserInfo 获取用户详情
func GetUserInfo(ctx *gin.Context) {
	tempUid, ok := ctx.Get("uid")
	uid := tempUid.(uint)
	if !ok {
		response.RespondWithStatusBadRequest(ctx, "获取用户信息失败")
		return
	}

	var userInfo *userModel.UserInfo
	var userCreationInfo *userModel.UserCreationInfo

	if temp, err := repository.GetUserInfo(uid); err != nil {
		response.RespondWithStatusBadRequest(ctx, "获取用户信息失败")
		return
	} else {
		userInfo = temp
		userInfo.AvatarUrl = utils.AddAvatarPrefix(userInfo.AvatarUrl)
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

	var userInfo *userModel.UserInfo
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
