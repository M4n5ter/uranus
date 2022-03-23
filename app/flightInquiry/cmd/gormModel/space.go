package gormModel

type Space struct {
	Model
	FlightInfoID uint
	//是否是头等舱
	IsFirstClass bool
	//价格
	Price uint64
	//剩余（由于有超卖可能性，可能为负）
	Surplus int64
}
