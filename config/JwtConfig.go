package config

import (
	"crypto/rand"
	"encoding/hex"
	"github.com/joho/godotenv"
	"note_app_server/global"
	"os"
)

func InitJWTConfig() {
	err := godotenv.Load()
	if err != nil {
		panic(".env file not found")
	}

	jwtKey := make([]byte, 32)
	jwtEnvParam, ok := os.LookupEnv("JWT_KEY")
	if !ok {
		if _, err = rand.Read(jwtKey); err != nil {
			panic("Failed to generate JWT Key")
		} else {
			err = godotenv.Write(map[string]string{"JWT_KEY": hex.EncodeToString(jwtKey)}, ".env")
			if err != nil {
				panic("Failed to write JWT Key")
			}
			global.JWTKey = jwtKey
		}
	} else {
		jwtKey, err = hex.DecodeString(jwtEnvParam)
		if err != nil {
			panic(err)
		}
		global.JWTKey = jwtKey
	}
}
