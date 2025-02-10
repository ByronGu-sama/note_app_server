package main

import (
	"note_app_server/config"
	"note_app_server/router"
	"note_app_server/service"
)

func main() {
	config.InitAppConfig()
	service.CreateJWTKey()
	//test.TestG()
	r := router.SetupRouter()
	if err := r.Run(config.AC.App.Port); err != nil {
		return
	}
}
