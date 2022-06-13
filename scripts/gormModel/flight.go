package gormModel

type Flight struct {
	Model
	//航班号
	Number string `gorm:"not null;default:unknown;unique;comment:航班号(YT1234)"`
	//机型
	FltType string `gorm:"not null;default:unknown;comment:机型"`
	//航班信息
	FlightInfos []FlightInfo `gorm:"foreignKey:FlightNumber;references:Number"`
}
