package config

import (
	"github.com/spf13/viper"
	"log"
	"note_app_server/model/appModel"
	"sync"
)

var (
	AC   *appModel.AppConfig
	once sync.Once
	wg   sync.WaitGroup
)

// InitAppConfig 读取config.yml文件
func InitAppConfig() {
	viper.SetConfigName("AppConfig")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config/yaml")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("read config failed, err:%v\n", err)
	}
	AC = &appModel.AppConfig{}
	if err := viper.Unmarshal(AC); err != nil {
		log.Fatalf("unmarshal config failed, err:%v\n", err)
	}

	once.Do(func() {
		wg.Add(7)
		go func() {
			defer wg.Done()
			InitMysqlConfig()
		}()
		go func() {
			defer wg.Done()
			InitOssConfig()
		}()
		go func() {
			defer wg.Done()
			InitCaptchaConfig()
		}()
		go func() {
			defer wg.Done()
			InitRedisConfig()
		}()
		go func() {
			defer wg.Done()
			InitElasticSearchConfig()
		}()
		go func() {
			defer wg.Done()
			InitMongoDB()
		}()
		go func() {
			defer wg.Done()
			InitJWTConfig()
		}()
	})
	wg.Wait()
}
