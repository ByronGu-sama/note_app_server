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
	"math/rand/v2"
	"mime/multipart"
	"net/http"
	"note_app_server/config"
	"note_app_server/global"
	"note_app_server/model"
	"note_app_server/repository"
	"note_app_server/response"
	"note_app_server/service"
	"note_app_server/utils"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
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

	noteId := utils.EncodeNoteId(fmt.Sprintf("%d-%d-%d", time.Now().Unix(), uid, rand.Int64()))

	var coverHeight int
	var cover string
	var picsNameList string
	var wg sync.WaitGroup
	type uploadRequest struct {
		fileHeader *multipart.FileHeader
		index      int
	}
	uploadChanList := make(chan uploadRequest)
	go func() {
		for req := range uploadChanList {
			wg.Add(1)
			go func(req uploadRequest) {
				defer wg.Done()
				openFile, err1 := req.fileHeader.Open()
				if err1 != nil {
					response.RespondWithStatusBadRequest(ctx, err1.Error())
					return
				}
				defer func(openFile multipart.File) {
					err := openFile.Close()
					if err != nil {
						log.Printf("close file err: %v", err)
					}
				}(openFile)

				// 判断文件类型
				fileType, err2 := utils.DetectFileType(&openFile)
				if err2 != nil {
					response.RespondWithStatusBadRequest(ctx, err2.Error())
					return
				}
				if fileType == "image/png" {
					fileType = "png"
				}
				if fileType == "image/jpeg" {
					fileType = "jpeg"
				}

				// 获取封面高度
				if req.index == 0 {
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
					coverHeight = img.Bounds().Dy()
				}

				fileName, err3 := service.UploadFileObject(config.AC.Oss.BucketName, "notePics/"+noteId+"/", openFile, fileType)
				// 获取封面
				if req.index == 0 {
					cover = fileName
				}

				if err3 != nil {
					response.RespondWithStatusInternalServerError(ctx, err3.Error())
					return
				} else {
					picsNameList += fileName + ";"
				}
			}(req)
		}
	}()
	for index, fileHeader := range ctx.Request.MultipartForm.File["file"] {
		uploadChanList <- uploadRequest{fileHeader, index}
	}
	close(uploadChanList)
	wg.Wait()

	note := &model.Note{
		Nid:         noteId,
		Uid:         uid,
		Cover:       cover,
		CoverHeight: coverHeight,
		Pics:        picsNameList[:len(picsNameList)-1],
		CategoryId:  1,
		Tags:        tags,
		Title:       title,
		Content:     content,
		Public:      1,
	}

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

	rawPicsList := strings.Split(note.Pics, ";")
	picsList := make([]string, 0)
	for _, pic := range rawPicsList {
		picsList = append(picsList, "http://"+config.AC.App.Host+config.AC.App.Port+"/note/pic/"+nid+"/"+pic)
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "success",
		"data": gin.H{
			"nid":              note.Nid,
			"uid":              note.Uid,
			"cover":            note.Cover,
			"pics":             picsList,
			"coverHeight":      note.CoverHeight,
			"title":            note.Title,
			"content":          note.Content,
			"createdAt":        note.CreatedAt,
			"updatedAt":        note.UpdatedAt,
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
	reader, err := service.GetOssObject(config.AC.Oss.BucketName, "notePics/"+nid+"/", fileName)

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

	for i := range result {
		result[i].AvatarUrl = "http://" + config.AC.App.Host + config.AC.App.Port + "/avatar/" + result[i].AvatarUrl
		result[i].Cover = "http://" + config.AC.App.Host + config.AC.App.Port + "/note/pic/" + result[i].Nid + "/" + result[i].Cover
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
		result[i].AvatarUrl = "http://" + config.AC.App.Host + config.AC.App.Port + "/avatar/" + result[i].AvatarUrl
		result[i].Cover = "http://" + config.AC.App.Host + config.AC.App.Port + "/note/pic/" + result[i].Nid + "/" + result[i].Cover
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
