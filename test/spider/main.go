package main

import (
	"context"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/chromedp/chromedp"
	"log"
	"strings"
	"time"
	"uranus/test/model"
)

//var url = "https://zh.flightaware.com/live/findflight?origin=ZBAA&destination=ZSHC"
var url = "https://flights.ctrip.com/schedule/hgh.sia.html?pageno=1"
var Flights = make([]model.Flight, 1)

func main() {
	options := []chromedp.ExecAllocatorOption{
		chromedp.Flag("headless", false),
		chromedp.Flag("blink-settings", "imagesEnabled=false"),
		chromedp.UserAgent(`Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/99.0.4844.51 Safari/537.36 Edg/99.0.1150.39`),
	}
	options = append(chromedp.DefaultExecAllocatorOptions[:], options...)
	ctx, _ := chromedp.NewExecAllocator(context.Background(), options...)
	chromeCtx, cancel := chromedp.NewContext(ctx, chromedp.WithLogf(log.Printf))
	defer cancel()
	err := chromedp.Run(chromeCtx, make([]chromedp.Action, 0, 1)...)
	if err != nil {
		log.Fatalln(err)
	}

	timeoutCtx1, cancel := context.WithTimeout(chromeCtx, time.Second*20)
	var flightNum string
	err = chromedp.Run(
		timeoutCtx1,
		chromedp.Navigate(url),
		chromedp.WaitVisible("#flt1 > tr:nth-child(1)"),
		chromedp.OuterHTML(`document.querySelector("body")`, &flightNum, chromedp.ByJSPath),
	)
	if err != nil {
		log.Fatalln(err)
	}
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(flightNum))
	if err != nil {
		log.Fatalln(err)
	}
	doc.Find(`#flt1 > tr`).Each(func(i int, selection *goquery.Selection) {
		if len(Flights) < i+1 {
			Flights = append(Flights, model.Flight{})
		}
		Flights[i].Number = selection.Find(`a`).Text()
	})

	timeoutCtx2, cancel := context.WithTimeout(chromeCtx, time.Second*20)
	var fltTypeJmp string
	err = chromedp.Run(
		timeoutCtx2,
		chromedp.Navigate(url),
		chromedp.WaitVisible("#flt1 > tr:nth-child(1) > td.flight_logo > span > span"),
		chromedp.OuterHTML(`document.querySelector("body")`, &fltTypeJmp, chromedp.ByJSPath),
	)
	doc, err = goquery.NewDocumentFromReader(strings.NewReader(fltTypeJmp))
	doc.Find(`#flt1 > tr`).Each(func(i int, selection *goquery.Selection) {
		if len(Flights) < i+1 {
			Flights = append(Flights, model.Flight{})
		}
		Flights[i].FltTypeJmp = selection.Find(`tr > td.flight_logo > span > span`).Text()
	})

	//timeoutCtx3, cancel := context.WithTimeout(chromeCtx, time.Second*20)
	//var SetOutDate string
	//err = chromedp.Run(
	//	timeoutCtx3,
	//	chromedp.Navigate(url),
	//	chromedp.WaitVisible("#flt1 > tr:nth-child(1) > td.price > div > input"),
	//	chromedp.OuterHTML(`document.querySelector("body")`, &SetOutDate, chromedp.ByJSPath),
	//)
	//doc, err = goquery.NewDocumentFromReader(strings.NewReader(SetOutDate))
	//doc.Find(`#flt1 > tr > td.price > div > input`).Each(func(i int, selection *goquery.Selection) {
	//	if len(Flights) < i+1 {
	//		Flights = append(Flights, model.Flight{})
	//	}
	//	Flights[i].FlightInfo = append(Flights[i].FlightInfo, )selection.Find(`tr > td.flight_logo > span > span`).Text()
	//})

	fmt.Println(len(Flights))
	for _, flight := range Flights {
		fmt.Printf("Number: %s, FltType: %s\n", flight.Number, flight.FltTypeJmp)
	}
}
