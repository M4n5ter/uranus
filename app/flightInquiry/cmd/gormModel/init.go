package gormModel

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
)

var GlobalDB *gorm.DB

func init() {
	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN:               "root:Happyboyhasfun1@(putpp.com:13306)/zero_gorm?charset=utf8mb4&parseTime=True&loc=Local",
		DefaultStringSize: 256,
	}), &gorm.Config{
		SkipDefaultTransaction:                   true,
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		log.Fatalln(err)
	}
	GlobalDB = db
}

type Model struct {
	gorm.Model
	DelState bool `gorm:"not null;default:0"`
	Version  int  `gorm:"not null;default:0"`
}
