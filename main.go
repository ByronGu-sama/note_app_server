package main

import (
	"note_app_server/config"
	"note_app_server/producer/connManager"
	"note_app_server/router"
)

func main() {
	// 初始化服务端配置
	config.InitAppConfig()
	// 初始化kafka连接
	connManager.InitKafkaConn()
	// 启动gin
	r := router.SetupRouter()
	if err := r.Run(config.AC.App.Port); err != nil {
		return
	}
}
