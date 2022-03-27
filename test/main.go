package main

import (
	"fmt"
	"github.com/jinzhu/copier"
)

type t1 struct {
	Name   string
	Number int64
}

type t2 struct {
	Name   string
	Number int64
}

func main() {
	test2 := make([]*t2, 3)
	test1 := []t1{
		{Name: "wyt1", Number: 66666},
		{Name: "wyt2", Number: 77777},
		{Name: "wyt3", Number: 88888},
	}
	_ = copier.Copy(&test2, &test1)
	for _, t := range test2 {
		fmt.Println(*t)
	}
	//gormModel2.GlobalDB.AutoMigrate(&gormModel2.FlightInfo{}, &gormModel2.Flight{}, &gormModel2.Space{}, &gormModel2.Ticket{}, &gormModel2.RefundAndChangeInfo{})
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
