package gormModel

type Space struct {
	Model
	FlightInfoID uint `gorm:"comment:对应的航班信息id"`
	//是否是头等舱
	IsFirstClass bool `gorm:"comment:是否是头等舱/商务舱"`
	//价格
	Price uint64 `gorm:"comment:价格"`
	//剩余（由于有超卖可能性，可能为负）
	Surplus int64 `gorm:"comment:剩余量"`
}
