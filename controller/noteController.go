package controller

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"log"
	"math"
	"math/rand/v2"
	"net/http"
	"note_app_server/config"
	"note_app_server/global"
	"note_app_server/model/noteModel"
	"note_app_server/model/userModel"
	"note_app_server/producer"
	"note_app_server/repository"
	"note_app_server/response"
	"note_app_server/service"
	"note_app_server/utils"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

// NewNote 创建笔记
func NewNote(ctx *gin.Context) {
	title := ctx.PostForm("title")
	content := ctx.PostForm("content")
	tags := ctx.PostForm("tags")
	if title == "" || content == "" {
		response.RespondWithStatusBadRequest(ctx, "关键信息不能为空")
		return
	}

	tempUid, ok := ctx.Get("uid")
	uid := tempUid.(uint)
	if !ok {
		response.RespondWithStatusBadRequest(ctx, "获取用户信息失败")
		return
	}

	noteId := utils.EncodeWithMD5(fmt.Sprintf("%d-%d-%d", time.Now().Unix(), uid, rand.Int64()))

	var coverHeight float64
	var cover string
	var picsNameList string
	fileList := make([]string, 0)

	for index, fileHeader := range ctx.Request.MultipartForm.File["file"] {
		openFile, err1 := fileHeader.Open()
		if err1 != nil {
			response.RespondWithStatusBadRequest(ctx, err1.Error())
			return
		}

		tempFile, err2 := io.ReadAll(openFile)
		if err2 != nil {
			response.RespondWithStatusBadRequest(ctx, err2.Error())
			return
		}
		// 重置指针
		_, seekErr := openFile.Seek(0, io.SeekStart)
		if seekErr != nil {
			response.RespondWithStatusBadRequest(ctx, seekErr.Error())
			return
		}

		// 判断文件类型
		fileType := utils.DetectFileType(tempFile)

		// 获取封面高度
		if index == 0 {
			img, _, err := image.Decode(openFile)
			if err != nil {
				response.RespondWithStatusBadRequest(ctx, err.Error())
				return
			}
			// 重置指针
			_, err = openFile.Seek(0, io.SeekStart)
			if err != nil {
				return
			}
			height := img.Bounds().Dy()
			width := img.Bounds().Dx()
			coverHeight = math.Round(100.0/float64(width)*float64(height)*100) / 100
		}

		fileName, err3 := service.UploadFileObject(config.AC.Oss.NotePicsBucket, noteId+"/", openFile, fileType)
		fileList = append(fileList, fileName)

		// 获取封面
		if index == 0 {
			cover = fileName
		}

		if err3 != nil {
			for i := range fileList {
				err := service.DeleteObject(config.AC.Oss.NotePicsBucket, noteId+"/", fileList[i])
				if err != nil {
					log.Println(err)
				}
			}
			response.RespondWithStatusInternalServerError(ctx, err3.Error())
			return
		} else {
			picsNameList += fileName + ";"
		}

		err := openFile.Close()
		if err != nil {
			log.Fatal(err)
		}
		openFile.Close()
	}

	curTime := time.Now()

	n := &noteModel.Note{
		Nid:         noteId,
		Uid:         uid,
		Cover:       cover,
		CoverHeight: coverHeight,
		Pics:        picsNameList[:len(picsNameList)-1],
		CategoryId:  1,
		Tags:        tags,
		Title:       title,
		CreatedAt:   curTime,
		UpdatedAt:   curTime,
		Content:     content,
		Public:      1,
	}

	tx := global.Db.Begin()
	if err := tx.Create(&n).Error; err != nil {
		tx.Rollback()
		response.RespondWithStatusBadRequest(ctx, "创建失败")
		return
	}
	userCreation := &userModel.UserCreationInfo{}
	if err := tx.Model(userCreation).Where("uid = ?", uid).Update("noteCount", gorm.Expr("noteCount + ?", 1)).Error; err != nil {
		tx.Rollback()
		response.RespondWithStatusBadRequest(ctx, "创建失败")
		return
	}
	tx.Commit()
	response.RespondWithStatusOK(ctx, "创建成功")

	tempUsername, _ := ctx.Get("username")
	username := tempUsername.(string)
	tempAvatarUrl, _ := ctx.Get("avatarUrl")
	avatarUrl := tempAvatarUrl.(string)

	esNote := &noteModel.ESNote{
		Nid:         noteId,
		Uid:         uid,
		Username:    username,
		AvatarUrl:   avatarUrl,
		Cover:       cover,
		CoverHeight: coverHeight,
		Pics:        picsNameList[:len(picsNameList)-1],
		Title:       title,
		Content:     content,
		LikesCount:  0,
		CreatedAt:   curTime,
		UpdatedAt:   curTime,
		Public:      true,
		CategoryId:  1,
		Tags:        tags,
		Status:      1,
	}

	producer.SyncToES(esNote)
}

// DelNote 删除笔记
func DelNote(ctx *gin.Context) {
	nid, uid, err := checkUidAndNid(ctx)
	if err != nil {
		response.RespondWithStatusBadRequest(ctx, err.Error())
		return
	}
	if err = producer.DelNote(uid, nid); err != nil {
		response.RespondWithStatusBadRequest(ctx, "删除失败")
		return
	}
	response.RespondWithStatusOK(ctx, "删除成功")
}

// EditNote 编辑笔记
func EditNote(ctx *gin.Context) {
	var note noteModel.Note
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

	rawPicsList := strings.Split(note.Pics, ";")
	picsList := make([]string, 0)
	for _, pic := range rawPicsList {
		picsList = append(picsList, utils.AddNotePicPrefix(nid, pic))
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "success",
		"data": gin.H{
			"nid":              note.Nid,
			"uid":              note.Uid,
			"avatarUrl":        utils.AddAvatarPrefix(note.AvatarUrl),
			"username":         note.Username,
			"pics":             picsList,
			"title":            note.Title,
			"content":          note.Content,
			"createdAt":        note.CreatedAt,
			"updatedAt":        note.UpdatedAt,
			"public":           note.Public,
			"categoryId":       note.CategoryId,
			"tags":             note.Tags,
			"likesCount":       note.LikesCount,
			"commentsCount":    note.CommentsCount,
			"collectionsCount": note.CollectionsCount,
			"sharesCount":      note.SharesCount,
			"viewsCount":       note.ViewsCount,
		},
	})
}

// GetNotePic 转换笔记图片地址
func GetNotePic(ctx *gin.Context) {
	nid := ctx.Param("nid")
	fileName := ctx.Param("fileName")
	reader, err := service.GetOssObject(config.AC.Oss.NotePicsBucket, nid+"/", fileName)

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
	fileType := filepath.Ext(fileName)
	if fileType == ".jpeg" || fileType == ".jpg" {
		fileType = "image/jpeg"
	}
	if fileType == ".png" {
		fileType = "image/png"
	}
	ctx.Header("Content-Type", fileType)
	ctx.Data(http.StatusOK, fileType, data)
}

// LikeNote 点赞笔记
func LikeNote(ctx *gin.Context) {
	nid, uid, err := checkUidAndNid(ctx)
	if err != nil {
		response.RespondWithStatusBadRequest(ctx, err.Error())
		return
	}

	if err = producer.LikeNote(uid, nid); err != nil {
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

	if err = producer.DislikeNote(uid, nid); err != nil {
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

	if err = producer.CollectNote(uid, nid); err != nil {
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

	if err = producer.AbandonNote(uid, nid); err != nil {
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

	for i := range result {
		result[i].AvatarUrl = utils.AddAvatarPrefix(result[i].AvatarUrl) + "?x-oss-process=style/compress_avatar"
		result[i].Cover = utils.AddNotePicPrefix(result[i].Nid, result[i].Cover) + "?x-oss-process=style/compress_cover"
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "success",
		"data":    result,
	})
}

// GetMyNotes 获取我的笔记列表
func GetMyNotes(ctx *gin.Context) {
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

	uid, ok := ctx.Get("uid")
	if !ok {
		response.RespondWithUnauthorized(ctx, "无权限")
		return
	}

	result, err3 := repository.GetNoteListWithUid(uid.(uint), page, limit)
	if err3 != nil {
		response.RespondWithStatusBadRequest(ctx, err3.Error())
		return
	}

	for i := range result {
		result[i].AvatarUrl = utils.AddAvatarPrefix(result[i].AvatarUrl) + "?x-oss-process=style/compress_avatar"
		result[i].Cover = utils.AddNotePicPrefix(result[i].Nid, result[i].Cover) + "?x-oss-process=style/compress_cover"
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "success",
		"data":    result,
	})
}

// GetNotesListWithKeyword 获取搜索帖子的结果
func GetNotesListWithKeyword(ctx *gin.Context) {
	keyword := ctx.Param("keyword")
	page := ctx.Query("page")
	limit := ctx.Query("limit")
	truePage, err1 := strconv.Atoi(page)
	if err1 != nil {
		response.RespondWithStatusBadRequest(ctx, err1.Error())
		return
	}
	trueLimit, err2 := strconv.Atoi(limit)
	if err2 != nil {
		response.RespondWithStatusBadRequest(ctx, err2.Error())
		return
	}
	if truePage <= 0 || trueLimit <= 0 {
		response.RespondWithStatusBadRequest(ctx, "参数错误")
		return
	}

	offset := (truePage - 1) * trueLimit

	if len(keyword) < 1 || len(keyword) > 200 {
		response.RespondWithStatusBadRequest(ctx, "关键词长度错误")
		return
	}
	result, err := repository.GetNoteListWithKeyword("notes", keyword, &offset, &trueLimit)
	if err != nil {
		response.RespondWithStatusBadRequest(ctx, err.Error())
		return
	}

	for i := range result {
		result[i].AvatarUrl = utils.AddAvatarPrefix(result[i].AvatarUrl) + "?x-oss-process=style/compress_avatar"
		result[i].Cover = utils.AddNotePicPrefix(result[i].Nid, result[i].Cover) + "?x-oss-process=style/compress_cover"
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
