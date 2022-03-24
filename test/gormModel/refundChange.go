package gormModel

import "time"

type RefundAndChangeInfo struct {
	Model
	//对应的票ID
	TicketID uint `gorm:"not null;default:0;comment:对应的票ID"`
	//是否退订(退订为真，改票为假)
	IsRefund bool `gorm:"not null;comment:1为退订，0为改票"`
	//时间1
	Time1 time.Time `gorm:"type:datetime(6);not null;comment:时间1"`
	//符合时间1时需要的费用(￥/人)
	Fee1 uint64 `gorm:"not null;default:99999;comment:符合时间1时需要的费用(￥/人)"`
	//时间2
	Time2 time.Time `gorm:"type:datetime(6);not null;comment:时间2"`
	//符合时间2时需要的费用(￥/人)
	Fee2 uint64 `gorm:"not null;default:99999;comment:符合时间2时需要的费用(￥/人)"`
	//时间3
	Time3 time.Time `gorm:"type:datetime(6);not null;comment:时间3"`
	//符合时间3时需要的费用(￥/人)
	Fee3 uint64 `gorm:"not null;default:99999;comment:符合时间3时需要的费用(￥/人)"`
	//时间4
	Time4 time.Time `gorm:"type:datetime(6);not null;comment:时间4"`
	//符合时间1时需要的费用(￥/人)
	Fee4 uint64 `gorm:"not null;default:99999;comment:符合时间4时需要的费用(￥/人)"`
	//时间5(可以为空，很多时候只有4个时间)
	Time5 time.Time `gorm:"type:datetime(6);comment:时间5"`
	//符合时间5时需要的费用(￥/人)
	Fee5 uint64 `gorm:"comment:符合时间5时需要的费用(￥/人)"`
}
