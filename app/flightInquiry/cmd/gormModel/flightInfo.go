package gormModel

import (
	"time"
)

// FlightInfo 航班信息
type FlightInfo struct {
	Model
	//对应的航班号
	FlightNumber string `gorm:"not null;default:unknown;comment:对应的航班号"`
	//出发日期
	SetOutDate time.Time `gorm:"type:datetime(6);not null;comment:出发日期"`
	//舱位
	Spaces Space
	//准点率
	Punctuality uint8 `gorm:"not null;default:0;comment:准点率(%)"`
	//起飞地点
	DepartPosition string `gorm:"not null;default:unknown;comment:起飞地点"`
	//起飞时间
	DepartTime time.Time `gorm:"type:datetime(6);not null;comment:起飞时间"`
	//降落地点
	ArrivePosition string `gorm:"not null;default:unknown;comment:降落地点"`
	//降落时间
	ArriveTime time.Time `gorm:"type:datetime(6);not null;comment:降落时间"`
}
