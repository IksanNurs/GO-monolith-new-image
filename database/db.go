package database

import (
	"log"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	db  *gorm.DB
	err error
)

func StartDB() {
	db, err = gorm.Open(mysql.Open(""+os.Getenv("DB_USER")+":"+os.Getenv("DB_PASSWORD")+"@tcp("+os.Getenv("DB_HOST")+")/"+os.Getenv("DB_NAME")+"?charset=utf8mb4&parseTime=True&loc=Local"), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	 db.Logger = db.Logger.LogMode(logger.Silent)
	// err = db.AutoMigrate(&models.User{}).Error
	// if err != nil {
	// 	log.Fatal(err)
	// }
}

func GetDB() *gorm.DB {
	return db
}
