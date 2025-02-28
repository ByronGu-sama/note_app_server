package main

import (
	"note_app_server/config"
	"note_app_server/router"
)

func main() {
	config.InitAppConfig()
	//test.TestRedis()
	r := router.SetupRouter()
	if err := r.Run(config.AC.App.Port); err != nil {
		return
	}
}
