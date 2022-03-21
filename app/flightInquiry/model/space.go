package model

import "gorm.io/gorm"

type Space struct {
	gorm.Model
	FlightInfoID uint
	//是否是头等舱
	IsFirstClass bool
	//价格
	Price uint64
}
