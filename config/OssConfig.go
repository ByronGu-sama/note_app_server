package config

import (
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"log"
	"note_app_server/global"
)

// InitOssConfig 初始化Oss配置
func InitOssConfig() {
	bucketName := AC.Oss.BucketName
	endPoint := AC.Oss.EndPoint
	region := AC.Oss.Region

	if endPoint == "" || bucketName == "" {
		log.Fatal("Please set yourEndpoint and bucketName.")
	}

	provider, err := oss.NewEnvironmentVariableCredentialsProvider()
	if err != nil {
		log.Fatal("new oss environment variable failed.")
	}

	// 创建OSSClient实例
	clientOptions := []oss.ClientOption{oss.SetCredentialsProvider(&provider)}
	// 配置bucket所在区域
	clientOptions = append(clientOptions, oss.Region(region))
	// 设置签名版本
	clientOptions = append(clientOptions, oss.AuthVersion(oss.AuthV4))
	client, err := oss.New(endPoint, "", "", clientOptions...)
	if err != nil {
		log.Fatal("new oss client failed.")
	}
	// 输出客户端信息。
	log.Printf("Oss客户端信息: %#v\n", client)
	global.OssClient = client
}
