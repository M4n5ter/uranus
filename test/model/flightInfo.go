package model

import (
	"gorm.io/gorm"
	"time"
)

// FlightInfo 航班信息
type FlightInfo struct {
	gorm.Model
	//对应的航班ID
	FlightID uint
	//出发日期
	SetOutDate time.Time
	//舱位
	Space Space
	//准点率
	Punctuality uint8
	//起飞地点
	StartPosition string
	//起飞时间
	StartTime time.Time
	//降落地点
	EndPosition string
	//降落时间
	EndTime time.Time
}
