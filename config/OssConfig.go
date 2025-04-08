package config

import (
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"log"
	"note_app_server/global"
	"sync"
)

// InitOssConfig 初始化Oss配置
func InitOssConfig() {
	avatarBucket := AC.Oss.AvatarBucket
	notePicsBucket := AC.Oss.NotePicsBucket
	endPoint := AC.Oss.EndPoint
	region := AC.Oss.Region

	if endPoint == "" || avatarBucket == "" || notePicsBucket == "" {
		log.Println("Please set yourEndpoint and bucketName.")
	}

	provider, err := oss.NewEnvironmentVariableCredentialsProvider()
	if err != nil {
		log.Println("new oss environment variable failed.")
	}

	var pool = &sync.Pool{
		New: func() interface{} {
			// 创建OSSClient实例
			clientOptions := []oss.ClientOption{oss.SetCredentialsProvider(&provider)}
			// 配置bucket所在区域
			clientOptions = append(clientOptions, oss.Region(region))
			// 设置签名版本
			clientOptions = append(clientOptions, oss.AuthVersion(oss.AuthV4))
			client, err := oss.New(endPoint, "", "", clientOptions...)
			if err != nil {
				log.Println("new oss client failed.")
			}
			return client
		},
	}
	global.OssClientPool = pool
}
