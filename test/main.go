package main

import "uranus/app/flightInquiry/cmd/gormModel"

func main() {
	gormModel.GlobalDB.AutoMigrate(&gormModel.FlightInfo{}, &gormModel.Flight{}, &gormModel.Space{}, &gormModel.Ticket{}, &gormModel.RefundAndChangeInfo{})
	//model.GlobalDB.AutoMigrate(&model.Space{})
	//s1 := model.Space{
	//	IsFirstClass: true,
	//	Price:        999,
	//}
	//fi1 := model.FlightInfo{
	//	Space:         s1,
	//	SetOutDate:    time.Date(2022, time.March, 22, 0, 0, 0, 0, time.Local),
	//	//	StartTime:     time.Now(),
	//	//	StartPosition: "你家",
	//	//	EndTime:       time.Now().Add(time.Hour * 5),
	//	//	EndPosition:   "我家",
	//	//	Punctuality:   97,
	//	//	FlightNumber:  "HG1099",
	//	//}
	//	f2 := model.Flight{
	//		FltTypeJmp: "320",
	//		Number:     "HG1099",
	//	}
	//	model.GlobalDB.Create(&f2)
	//	//model.GlobalDB.Create(&fi1)
	//	//model.GlobalDB.First(&model.Flight{}).Where("number = ?", fi1.FlightNumber).Association("FlightInfos").Append(&fi1)
}
