package config

import (
	captcha20230305 "github.com/alibabacloud-go/captcha-20230305/client"
	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/aliyun/credentials-go/credentials"
	"log"
	"note_app_server/global"
	"os"
	"sync"
)

func InitCaptchaConfig() {
	config := new(credentials.Config).
		SetType("access_key").
		SetAccessKeyId(os.Getenv("OSS_ACCESS_KEY_ID")).
		SetAccessKeySecret(os.Getenv("OSS_ACCESS_KEY_SECRET"))

	akCredential, err := credentials.NewCredential(config)
	if err != nil {
		panic(err)
	}
	accessKeyId, err1 := akCredential.GetAccessKeyId()
	if err1 != nil {
		panic(err1)
	}
	accessSecret, err2 := akCredential.GetAccessKeySecret()
	if err2 != nil {
		panic(err2)
	}

	cfg := &openapi.Config{}
	cfg.AccessKeyId = accessKeyId
	cfg.AccessKeySecret = accessSecret

	pool := &sync.Pool{
		New: func() interface{} {
			cfg.Endpoint = tea.String(AC.Captcha.EndPoint)
			cfg.ConnectTimeout = tea.Int(5000)
			cfg.ReadTimeout = tea.Int(5000)
			client, err3 := captcha20230305.NewClient(cfg)
			if err3 != nil {
				log.Fatal(err3)
			}
			return client
		},
	}

	global.CaptchaClientPool = pool
}
