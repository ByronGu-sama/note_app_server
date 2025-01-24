package main

import (
	"note_app_server1/config"
	"note_app_server1/router"
	"note_app_server1/test"
)

func main() {
	config.InitAppConfig()
	test.Test()

	r := router.SetupRouter()
	if err := r.Run(config.AC.App.Port); err != nil {
		return
	}
}
