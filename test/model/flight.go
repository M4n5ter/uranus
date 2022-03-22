package model

import (
	"gorm.io/gorm"
)

type Flight struct {
	gorm.Model
	//航班号
	Number string `gorm:"unique"`
	//机型
	FltTypeJmp string
	//航班信息
	FlightInfos []FlightInfo `gorm:"foreignKey:FlightNumber;references:Number"`
}
