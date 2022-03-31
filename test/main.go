package main

import (
	"fmt"
	"github.com/pkg/errors"
	"uranus/common/xerr"
)

func main() {
	e := errors.Wrapf(xerr.NewErrMsg("不支持更改此状态,不支持更改为待支付状态"), "hi")
	err := errors.WithMessagef(e, "hello")
	fmt.Println(e)
	fmt.Println(err)
	//t, err := time.Parse("2006-01-02", "2006-01-02 00:00:00")
	//fmt.Println(t.Format("2006-01-02 15:04:05"), err)
	//gormModel2.GlobalDB.AutoMigrate(&gormModel2.FlightInfo{}, &gormModel2.Flight{}, &gormModel2.Space{}, &gormModel2.Ticket{}, &gormModel2.RefundAndChangeInfo{})
	//commonModel.GlobalDB.AutoMigrate(&commonModel.Space{})
	//s1 := commonModel.Space{
	//	IsFirstClass: true,
	//	Price:        999,
	//}
	//fi1 := commonModel.FlightInfo{
	//	Space:         s1,
	//	SetOutDate:    time.Date(2022, time.March, 22, 0, 0, 0, 0, time.Local),
	//	//	StartTime:     time.Now(),
	//	//	StartPosition: "你家",
	//	//	EndTime:       time.Now().Add(time.Hour * 5),
	//	//	EndPosition:   "我家",
	//	//	Punctuality:   97,
	//	//	FlightNumber:  "HG1099",
	//	//}
	//	f2 := commonModel.Flight{
	//		FltTypeJmp: "320",
	//		Number:     "HG1099",
	//	}
	//	commonModel.GlobalDB.Create(&f2)
	//	//commonModel.GlobalDB.Create(&fi1)
	//	//commonModel.GlobalDB.First(&commonModel.Flight{}).Where("number = ?", fi1.FlightNumber).Association("FlightInfos").Append(&fi1)
}
