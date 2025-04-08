package utils

import (
	"log"
)

func SafeGo(fn func()) {
	defer func() {
		if err := recover(); err != nil {
			log.Println("Recovered from panic")
		}
	}()
	fn()
}
