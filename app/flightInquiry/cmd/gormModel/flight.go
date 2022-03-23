package gormModel

type Flight struct {
	Model
	//航班号
	Number string `gorm:"unique;comment:航班号"`
	//机型
	FltTypeJmp string `gorm:"comment:机型"`
	//航班信息
	FlightInfos []FlightInfo `gorm:"foreignKey:FlightNumber;references:Number"`
}
