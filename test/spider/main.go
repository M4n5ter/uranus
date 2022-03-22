package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/chromedp/chromedp"
	"log"
	"net/http"
	"runtime"
	"strconv"
	"strings"
	"time"
	"uranus/test/model"
)

var CNFlightNames = `北京首都国际机场 PEK 北京
上海浦东国际机场 PVG 上海
广州白云国际机场 CAN 广东广州
上海虹桥国际机场 SHA 上海
深圳宝安国际机场 SZX 广东深圳
成都双流国际机场 CTU 四川成都
昆明巫家坝国际机场 KMG 云南昆明
海口美兰机场 HAK 海南海口
西安咸阳国际机场 SIA 陕西西安
杭州萧山国际机场 HGH 浙江杭州
厦门高崎国际机场 XMN 福建厦门
重庆江北国际机场 CKG 重庆
青岛流亭机场 TAO 山东青岛
大连周水子国际机场 DLC 辽宁大连
南京禄口国际机场 NKG 江苏南京
武汉天河机场 WUH 湖北武汉
沈阳桃仙国际机场 SHE 辽宁沈阳
乌鲁木齐地窝堡国际机场 URC 新疆乌鲁木齐
长沙黄花国际机场 CSX 湖北长沙
福州长乐国际机场 FOC 福建福州
桂林两江机场 KWL 广西桂林
哈尔滨太平国际机场 HRB 黑龙江哈尔滨
贵阳龙洞堡机场 KWE 贵州贵阳
郑州新郑国际机场 CGO 河南郑州
三亚凤凰机场 SYX 海南三亚
温州永强机场 WNZ 浙江温州
济南遥墙机场 TNA 山东济南
宁波栎社机场 NGB 浙江宁波
天津滨海国际机场 TSN 天津
太原武宿机场 TYN 山西太原
南宁吴圩机场 NNG 广西南宁
南昌昌北机场 KHN 江西南昌
长春大房身机场 CGQ 吉林长春
张家界荷花机场 DYG 湖南张家界
合肥骆岗机场 HFE 安徽合肥
西双版纳嘎洒机场 JHG 云南西双版纳
泉州晋江机场 JJN 福建晋江
兰州中川机场 LHW 甘肃兰州
烟台莱山机场 YNT 山东烟台
九寨黄龙机场 JZH 四川九寨沟
丽江三义机场 LJG 云南丽江
汕头外砂机场 SWA 广东汕头
呼和浩特据白塔机场 HET 内蒙古呼和浩特
拉萨贡嘎机场 LXA 西藏拉萨
珠海三灶机场 ZUH 广东珠海
银川河东机场 INC 宁夏银川
延吉朝阳川机场 YNJ 吉林延吉
武夷山机场 WUS 福建武夷山
西宁曹家堡机场 XNN 青海西宁
湛江机场 ZHA 广东湛江
舟山机场 HSN 浙江舟山
黄山屯溪机场 TXN 安徽黄山
宜昌三峡机场 YIH 湖北宜昌
喀什机场 KHG 新疆喀什
无锡机场 WUX 江苏无锡
包头二里半机场 BAV 内蒙古包头
伊宁机场 YIN 新疆伊宁
大理机场 DLU 云南大理
北海福成机场 BHY 广西北海
石家庄正定机场 SJW 河北石家庄
常州奔牛机场 CZX 江苏常州
库尔勒机场 KRL 新疆库尔勒`

//var url = "https://zh.flightaware.com/live/findflight?origin=ZBAA&destination=ZSHC"
const url = "https://flights.ctrip.com/schedule/$1.$2.html?pageno="

var cancel context.CancelFunc
var pool = make(chan struct{}, 3)

func main() {
	defer cancel()
	list := strings.Split(CNFlightNames, "\n")
	flightMap := make(map[string]string)
	for _, s := range list {
		flightMap[strings.Split(s, " ")[2]] = strings.ToLower(strings.Split(s, " ")[1])
	}
	for _, v1 := range flightMap {
		tempUrl := strings.Replace(url, "$1", v1, -1)
		pool <- struct{}{}
		go func() {
			for _, v2 := range flightMap {
				tempUrl1 := strings.Replace(tempUrl, "$2", v2, -1)
				i := 1
				for {
					tempUrl2 := tempUrl1 + strconv.Itoa(i)
					resp, err := http.Get(tempUrl2)
					if err != nil {
						break
					}
					if resp.StatusCode == http.StatusOK {
						cancel, err = spider(tempUrl2)
						if err != nil {
							//time.Sleep(time.Second)
							cancel()
							break
						}
						cancel()
						i++
					}
				}
			}
			<-pool
		}()
	}

}

func spider(url string) (cancel context.CancelFunc, err error) {
	var Flights = make([]model.Flight, 1)
	var FlightInfos = make([]model.FlightInfo, 1)
	options := []chromedp.ExecAllocatorOption{
		chromedp.Flag("headless", false),
		chromedp.Flag("blink-settings", "imagesEnabled=false"),
		chromedp.UserAgent(`Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/99.0.4844.51 Safari/537.36 Edg/99.0.1150.39`),
	}
	options = append(chromedp.DefaultExecAllocatorOptions[:], options...)
	ctx, _ := chromedp.NewExecAllocator(context.Background(), options...)
	chromeCtx, cancel := chromedp.NewContext(ctx, chromedp.WithLogf(log.Printf))
	defer cancel()
	err = chromedp.Run(chromeCtx, make([]chromedp.Action, 0, 1)...)
	if err != nil {
		log.Fatalln(err)
	}

	timeoutCtx1, cancel := context.WithTimeout(chromeCtx, time.Hour*20)
	var htmlBody string
	err = chromedp.Run(
		timeoutCtx1,
		chromedp.Navigate(url),
		chromedp.OuterHTML(`document.querySelector("body")`, &htmlBody, chromedp.ByJSPath),
	)
	if err != nil {
		log.Fatalln(err)
	}
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlBody))
	if err != nil {
		log.Fatalln(err)
	}
	if strings.Contains(doc.Find(`#flt1 > tr > td > div > div > div > div`).Text(), "无直飞航班") {
		fmt.Println("无直飞航班")
		return cancel, errors.New("无直飞航班")
	}
	if strings.Contains(doc.Find(`#search_schedule > div > div > div.recommend_box > div.recommend_t > h3`).Text(), "热门航班推荐") {
		fmt.Println("意外的页面")
		return cancel, errors.New("意外的页面")
	}
	doc.Find(`#flt1 > tr`).Each(func(i int, selection *goquery.Selection) {
		//爬flight
		if len(Flights) < i+1 {
			Flights = append(Flights, model.Flight{})
		}
		Flights[i].FltTypeJmp = selection.Find(`tr > td.flight_logo > span > span`).Text()
		num := selection.Find(`a`).Text()
		Flights[i].Number = num
		//爬flightInfo
		if len(FlightInfos) < i+1 {
			FlightInfos = append(FlightInfos, model.FlightInfo{})
		}
		FlightInfos[i].FlightNumber = num
		p, _ := strconv.Atoi(strings.TrimRight(selection.Find(`td.punctuality`).Text(), "%"))
		FlightInfos[i].Punctuality = uint8(p)
		FlightInfos[i].StartPosition = selection.Find(`td.depart > div`).Text()
		FlightInfos[i].EndPosition = selection.Find(`td.arrive > div`).Text()
		setOutTime, _ := selection.Find(`td.price > div > input`).Attr("value")
		sod, err := time.ParseInLocation("2006-01-02", setOutTime, time.Local)
		if err != nil {
			return
		}
		FlightInfos[i].SetOutDate = sod
		departTime := selection.Find(`td.depart > strong`).Text()
		arriveTime := selection.Find(`td.arrive > strong`).Text()
		fullDepartTime := setOutTime + " " + departTime + ":00"
		fullArriveTime := setOutTime + " " + arriveTime + ":00"
		dt, err := time.ParseInLocation("2006-01-02 15:04:05", fullDepartTime, time.Local)
		if err != nil {
			return
		}
		at, err := time.ParseInLocation("2006-01-02 15:04:05", fullArriveTime, time.Local)
		if err != nil {
			return
		}
		FlightInfos[i].StartTime, FlightInfos[i].EndTime = dt, at
	})

	//timeoutCtx2, cancel := context.WithTimeout(chromeCtx, time.Hour*20)
	//var fltTypeJmp string
	//err = chromedp.Run(
	//	timeoutCtx2,
	//	chromedp.Navigate(url),
	//	chromedp.WaitVisible("#flt1 > tr:nth-child(1) > td.flight_logo > span > span"),
	//	chromedp.OuterHTML(`document.querySelector("body")`, &fltTypeJmp, chromedp.ByJSPath),
	//)
	//doc, err = goquery.NewDocumentFromReader(strings.NewReader(fltTypeJmp))
	//doc.Find(`#flt1 > tr`).Each(func(i int, selection *goquery.Selection) {
	//	if len(Flights) < i+1 {
	//		Flights = append(Flights, model.Flight{})
	//	}
	//	Flights[i].FltTypeJmp = selection.Find(`tr > td.flight_logo > span > span`).Text()
	//})
	//
	//fmt.Println(len(Flights))
	//for _, flight := range Flights {
	//	if len(flight.Number) != 0 && len(flight.FltTypeJmp) != 0 {
	//		model.GlobalDB.Create(&model.Flight{Number: flight.Number, FltTypeJmp: flight.FltTypeJmp})
	//	}
	//}

	for _, info := range FlightInfos {
		fmt.Println(
			info.FlightNumber,
			info.Punctuality,
			info.StartPosition,
			info.StartTime,
			info.EndPosition,
			info.EndTime,
			info.SetOutDate)
	}
	model.GlobalDB.Create(&Flights)
	model.GlobalDB.Create(&FlightInfos)
	runtime.GC()
	return
}
