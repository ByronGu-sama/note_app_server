package config

import (
	"github.com/spf13/viper"
	"log"
	"note_app_server1/model"
	"sync"
)

var AC *model.AppConfig

// InitAppConfig 读取config.yml文件
func InitAppConfig() {
	viper.SetConfigName("AppConfig")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config/yaml")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("read config failed, err:%v\n", err)
	}
	AC = &model.AppConfig{}
	if err := viper.Unmarshal(AC); err != nil {
		log.Fatalf("unmarshal config failed, err:%v\n", err)
	}

	var once sync.Once
	var wg sync.WaitGroup
	once.Do(func() {
		wg.Add(3)
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
			InitRedisConfig()
		}()
	})
	wg.Wait()
}
