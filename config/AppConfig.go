package config

import (
	"github.com/spf13/viper"
	"log"
)

type AppConfig struct {
	App struct {
		Name string
		Port string
	}
	Mysql struct {
		Dsn             string
		MaxIdleConns    int
		MaxOpenConns    int
		ConnMaxLifetime string
	}
}

var AC *AppConfig

// InitAppConfig 读取config.yml文件
func InitAppConfig() {
	viper.SetConfigName("AppConfig")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config/yaml")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("read config failed, err:%v\n", err)
	}
	AC = &AppConfig{}
	if err := viper.Unmarshal(AC); err != nil {
		log.Fatalf("unmarshal config failed, err:%v\n", err)
	}
	InitMysqlConfig()
}
