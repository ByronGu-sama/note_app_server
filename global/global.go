package global

import (
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"gorm.io/gorm"
)

var (
	Db        *gorm.DB
	OssClient *oss.Client
)
