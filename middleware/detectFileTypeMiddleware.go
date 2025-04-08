package middleware

import (
	"github.com/gabriel-vasile/mimetype"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"note_app_server/response"
	"note_app_server/utils"
	"sync"
)

// DetectNotePicsTypeMiddleware 检测笔记图片格式
func DetectNotePicsTypeMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		err := c.Request.ParseMultipartForm(100 << 20)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Unable to parse form"})
			c.Abort()
			return
		}
		var wg sync.WaitGroup
		for _, fileHeader := range c.Request.MultipartForm.File["file"] {
			wg.Add(1)
			go func() {
				utils.SafeGo(func() {
					defer wg.Done()
					openFile, err1 := fileHeader.Open()
					defer openFile.Close()

					if err1 != nil {
						response.RespondWithStatusBadRequest(c, err1.Error())
						c.Abort()
						return
					}

					byteData, err := io.ReadAll(openFile)
					if err != nil {
						response.RespondWithStatusBadRequest(c, err.Error())
						c.Abort()
						return
					}
					if _, err := (openFile).Seek(0, io.SeekStart); err != nil {
						response.RespondWithStatusBadRequest(c, err.Error())
						c.Abort()
						return
					}

					fileType := mimetype.Detect(byteData)

					if fileType.String() != "image/png" && fileType.String() != "image/jpeg" && fileType.String() != "image/webp" && fileType.String() != "image/heic" && fileType.String() != "image/heif" {
						response.RespondWithStatusBadRequest(c, "不支持的文件格式:"+fileType.String())
						c.Abort()
						return
					}
				})
			}()
		}
		wg.Wait()
		c.Next()
	}
}

// DetectNormalImageTypeMiddleware 检测文件格式
func DetectNormalImageTypeMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		avatar, err := c.FormFile("file")
		if err != nil {
			response.RespondWithStatusBadRequest(c, err.Error())
			c.Abort()
			return
		}

		openFile, err1 := avatar.Open()
		if err1 != nil {
			response.RespondWithStatusBadRequest(c, err1.Error())
			c.Abort()
			return
		}
		defer openFile.Close()

		byteData, err := io.ReadAll(openFile)
		if _, err := (openFile).Seek(0, io.SeekStart); err != nil {
			response.RespondWithStatusBadRequest(c, err.Error())
			c.Abort()
			return
		}

		if err != nil {
			response.RespondWithStatusBadRequest(c, err.Error())
			c.Abort()
			return
		}

		fileType := mimetype.Detect(byteData)

		if fileType.String() != "image/png" && fileType.String() != "image/jpeg" && fileType.String() != "image/webp" && fileType.String() != "image/heic" && fileType.String() != "image/heif" {
			response.RespondWithStatusBadRequest(c, "不支持的文件格式:"+fileType.String())
			c.Abort()
			return
		}

		c.Next()
	}
}
