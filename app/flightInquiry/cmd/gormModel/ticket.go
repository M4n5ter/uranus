package gormModel

// Ticket 票
type Ticket struct {
	Model
	//对应舱位ID
	SpaceID uint `gorm:"not null;default:0;comment:对应舱位ID"`
	//价格
	Price uint64 `gorm:"not null;default:999999;comment:价格(￥)"`
	//折扣(15 表示八五折)
	Discount uint8 `gorm:"not null;default:0;comment:折扣(-n%)"`
	//托运行李额(Checked baggage allowance)(KG)
	CBA uint8 `gorm:"not null;default:20;comment:托运行李额(KG)"`
	//退改票信息
	RefundAndChangeInfos []RefundAndChangeInfo
}
