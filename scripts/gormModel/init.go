package gormModel

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"time"
)

var GlobalDB *gorm.DB

func init() {
	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN:               "root:flightpassword@(putpp.com:13306)/scripts?charset=utf8mb4&parseTime=True&loc=Local",
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
	ID         uint      `gorm:"primarykey;"`
	CreateTime time.Time `gorm:"type:datetime(6);not null;"`
	UpdateTime time.Time `gorm:"type:datetime(6);not null;"`
	DeleteTime time.Time `gorm:"type:datetime(6);index;;not null;"`
	DelState   bool      `gorm:"not null;default:0;comment:是否已经删除"`
	Version    int       `gorm:"not null;default:0;comment:版本号"`
}
