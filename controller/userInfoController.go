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
	"time"
)

// GetAvatarUrl 获取代理头像地址
func GetAvatarUrl(ctx *gin.Context) {
	fileName := ctx.Param("fileName")
	reader, err := service.GetOssObject(config.AC.Oss.AvatarBucket, "", fileName)

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
			"gender":    userInfo.Gender,
			"birth":     userInfo.Birth.Format("2006-01-02"),
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
	username := ctx.PostForm("username")
	signature := ctx.PostForm("signature")
	birth := ctx.PostForm("birth")
	gender := ctx.PostForm("gender")
	avatarFile, err := ctx.FormFile("avatarFile")
	hasNewAvatar := false
	if err == nil {
		hasNewAvatar = true
	}

	tempUid, ok := ctx.Get("uid")
	uid := tempUid.(uint)
	if !ok {
		response.RespondWithStatusBadRequest(ctx, "获取用户信息失败")
		return
	}
	if username == "" || gender == "" {
		response.RespondWithStatusBadRequest(ctx, "关键信息不能为空")
		return
	}
	if gender != "0" && gender != "1" && gender != "2" {
		response.RespondWithStatusBadRequest(ctx, "性别信息错误")
		return
	}

	birthTime, err1 := time.Parse("2006-01-02", birth)
	if err1 != nil {
		response.RespondWithStatusBadRequest(ctx, err1.Error())
		return
	}

	oldAvatar, newAvatar := "", ""
	if hasNewAvatar {
		// 转换为file类型的文件
		openFile, err2 := avatarFile.Open()
		if err2 != nil {
			response.RespondWithStatusBadRequest(ctx, err2.Error())
			return
		}

		// 检查文件类型
		contentType, err3 := utils.DetectFileType(&openFile)
		if err3 != nil {
			response.RespondWithStatusBadRequest(ctx, err3.Error())
			return
		}
		if contentType != "image/jpeg" && contentType != "image/png" {
			response.RespondWithStatusBadRequest(ctx, "不支持的文件类型")
			return
		}
		if contentType == "image/png" {
			contentType = "png"
		}
		if contentType == "image/jpeg" {
			contentType = "jpeg"
		}

		oldAvatar, err = repository.GetLastAvatarUrl(uid)
		if err != nil {
			log.Fatal(err)
		}

		// 上传文件至oss
		newAvatar, err = service.UploadFileObject(config.AC.Oss.AvatarBucket, "", openFile, contentType)
		if err != nil {
			response.RespondWithStatusBadRequest(ctx, err.Error())
			return
		}
	}

	userInfo := &userModel.UserInfo{}
	if hasNewAvatar {
		userInfo.Uid = uid
		userInfo.Username = username
		userInfo.AvatarUrl = newAvatar
		userInfo.Birth = birthTime
		userInfo.Gender = gender
		userInfo.Signature = signature
	} else {
		userInfo.Uid = uid
		userInfo.Username = username
		userInfo.Birth = birthTime
		userInfo.Gender = gender
		userInfo.Signature = signature
	}

	if err6 := repository.UpdateUserInfo(userInfo); err6 != nil {
		if hasNewAvatar {
			if err7 := service.DeleteObject(config.AC.Oss.AvatarBucket, "", newAvatar); err7 != nil {
				log.Fatal(err7)
			}
		}
		response.RespondWithStatusBadRequest(ctx, "更新失败")
		return
	}
	if hasNewAvatar {
		if err8 := service.DeleteObject(config.AC.Oss.AvatarBucket, "", oldAvatar); err8 != nil {
			log.Fatal(err8)
		}
	}

	response.RespondWithStatusOK(ctx, "更新成功")
}
