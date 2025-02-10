package config

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"note_app_server/global"
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
	duration, err := time.ParseDuration(AC.Mysql.ConnMaxLifetime)
	if err != nil {
		sqlDB.SetConnMaxLifetime(time.Hour)
	}
	sqlDB.SetConnMaxLifetime(duration)
	global.Db = db
}
