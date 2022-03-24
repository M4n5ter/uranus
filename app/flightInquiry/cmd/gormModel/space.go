package gormModel

// Space 舱位
type Space struct {
	Model
	FlightInfoID uint `gorm:"not null;default:0;comment:对应的航班信息id"`
	//是否是头等舱
	IsFirstClass bool `gorm:"not null;default:0;comment:是否是头等舱/商务舱"`
	//总量
	Total int64 `gorm:"not null;default:0;comment:总量"`
	//剩余（由于有超卖可能性，可能为负）
	Surplus int64 `gorm:"not null;default:0;comment:剩余量"`
	//票
	Tickets []Ticket
}
