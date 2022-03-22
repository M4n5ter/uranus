package model

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
)

var GlobalDB *gorm.DB

func init() {
	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN:               "root:Happyboyhasfun1@tcp(putpp.com:13306)/gorm_test?charset=utf8mb4&parseTime=True&loc=Local",
		DefaultStringSize: 256,
	}), &gorm.Config{
		SkipDefaultTransaction:                   true,
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		log.Fatalln(err)
	}
	GlobalDB = db
	GlobalDB.AutoMigrate(&Flight{}, &FlightInfo{}, &Space{})
}
