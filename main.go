package main

import (
	"note_app_server1/config"
	"note_app_server1/router"
)

func main() {
	config.InitAppConfig()
	r := router.SetupRouter()
	if err := r.Run(config.AC.App.Port); err != nil {
		return
	}
}
