package config

import (
	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	green20220302 "github.com/alibabacloud-go/green-20220302/v2/client"
	"github.com/alibabacloud-go/tea/tea"
	"note_app_server/global"
	"os"
)

func InitContentCheckConfig() {
	tempKeyId := os.Getenv("OSS_ACCESS_KEY_ID")
	tempKeySecret := os.Getenv("OSS_ACCESS_KEY_SECRET")

	config := &openapi.Config{
		AccessKeyId:     &tempKeyId,
		AccessKeySecret: &tempKeySecret,
		RegionId:        tea.String("cn-shanghai"),
		Endpoint:        tea.String("green-cip.cn-shanghai.aliyuncs.com"),
		/**
		 * 请设置超时时间。服务端全链路处理超时时间为10秒，请做相应设置。
		 * 如果您设置的ReadTimeout小于服务端处理的时间，程序中会获得一个ReadTimeout异常。
		 */
		ConnectTimeout: tea.Int(3000),
		ReadTimeout:    tea.Int(6000),
	}
	client, err := green20220302.NewClient(config)
	if err != nil {
		panic(err)
	}
	global.ContentCheckClient = client
}
