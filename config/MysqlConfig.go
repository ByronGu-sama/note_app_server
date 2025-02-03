package config

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"note_app_server/global"
	"note_app_server/utils"
	"time"
)

func InitMysqlConfig() {
	dsn := AC.Mysql.Dsn
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal(err)
	}
	sqlDB.SetMaxIdleConns(AC.Mysql.MaxIdleConns)
	sqlDB.SetMaxOpenConns(AC.Mysql.MaxOpenConns)
	duration, err := utils.AtoT(AC.Mysql.ConnMaxLifetime)
	if err != nil {
		sqlDB.SetConnMaxLifetime(time.Hour)
	}
	sqlDB.SetConnMaxLifetime(duration)
	global.Db = db
}
