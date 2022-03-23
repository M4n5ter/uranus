package gormModel

import (
	"time"
)

// FlightInfo 航班信息
type FlightInfo struct {
	Model
	//对应的航班号
	FlightNumber string `gorm:"comment:对应的航班号"`
	//出发日期
	SetOutDate time.Time `gorm:"comment:出发日期"`
	//舱位
	Space Space
	//准点率
	Punctuality uint8 `gorm:"comment:准点率"`
	//起飞地点
	StartPosition string `gorm:"comment:起飞地点"`
	//起飞时间
	StartTime time.Time `gorm:"comment:起飞时间"`
	//降落地点
	EndPosition string `gorm:"comment:降落地点"`
	//降落时间
	EndTime time.Time `gorm:"comment:降落时间"`
}
