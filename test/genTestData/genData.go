package main

import (
	"github.com/panjf2000/ants/v2"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"math/rand"
	"sync"
	"time"
	"uranus/commonModel"
)

type Config struct {
	DB struct {
		DataSource string
	}
	Cache cache.CacheConf
}

func main() {
	var c Config
	conf.MustLoad("etc/genData.yaml", &c)
	insertSpaces(c, 1, 18581)
	insertTickets(c, 1, 37162)
	insertRCIs(c, 1, 37162)
}

func RandInt(min, max int, id int64) int {
	if min == max {
		return min
	}
	if max < min {
		min, max = max, min
	}
	r := rand.New(rand.NewSource(time.Now().UnixNano() + id))
	return r.Intn(max-min) + min
}

// 生成一个 flightInfoID 对应的一对 space
func genACoupleRandomSpaces(flightInfoID int64) []*commonModel.Spaces {
	space0 := &commonModel.Spaces{
		FlightInfoId: flightInfoID,
		IsFirstClass: 0,
		Total:        int64(RandInt(50, 100, flightInfoID)),
		Surplus:      int64(RandInt(50, 100, flightInfoID)),
		LockedStock:  0,
	}

	space1 := &commonModel.Spaces{
		FlightInfoId: flightInfoID,
		IsFirstClass: 1,
		Total:        int64(RandInt(50, 100, flightInfoID)),
		Surplus:      int64(RandInt(50, 100, flightInfoID)),
		LockedStock:  0,
	}

	return []*commonModel.Spaces{space0, space1}
}

func insertSpaces(c Config, minFlightInfoID, maxFlightInfoID int64) {
	var wg sync.WaitGroup
	spacesModel := commonModel.NewSpacesModel(sqlx.NewMysql(c.DB.DataSource), c.Cache)
	pool, _ := ants.NewPoolWithFunc(1000, func(flightInfoId interface{}) {
		twoSpaces := genACoupleRandomSpaces(flightInfoId.(int64))
		_, _ = spacesModel.Insert(nil, twoSpaces[0])
		_, _ = spacesModel.Insert(nil, twoSpaces[1])
		wg.Done()
	})
	defer pool.Release()

	for i := minFlightInfoID; i <= maxFlightInfoID; i++ {
		wg.Add(1)
		_ = pool.Invoke(i)
	}

	wg.Wait()
}

func genRandomTicket(spaceID int64) *commonModel.Tickets {
	return &commonModel.Tickets{
		SpaceId:  spaceID,
		Price:    int64(RandInt(50000, 500000, spaceID)),
		Discount: int64(RandInt(0, 90, spaceID)),
		Cba:      20,
	}
}

func insertTickets(c Config, minSpaceID, maxSpaceID int64) {
	var wg sync.WaitGroup
	ticketsModel := commonModel.NewTicketsModel(sqlx.NewMysql(c.DB.DataSource), c.Cache)
	pool, _ := ants.NewPoolWithFunc(50000, func(spaceId interface{}) {
		_, err := ticketsModel.Insert(nil, genRandomTicket(spaceId.(int64)))
		if err != nil {
			logx.Error(err)
		}

		wg.Done()
	})
	defer pool.Release()

	for i := minSpaceID; i <= maxSpaceID; i++ {
		wg.Add(1)
		_ = pool.Invoke(i)
	}

	wg.Wait()
}

// 因旅客原因，申请退票，在航班规定起飞时间：24小时以前，应支付原票款10%的退票费。
// 24小时以内，2小时以前，应支付原票款20%的退票费;2小时以内，应支付原票款30%的退票费。
// 误机，应支付原票款50%的退票费。todo
func genACoupleOfRandomRCI(ticketId int64) []*commonModel.RefundAndChangeInfos {
	rci0 := &commonModel.RefundAndChangeInfos{
		TicketId: ticketId,
		IsRefund: 0,
		Time1:    time.Now().AddDate(0, 0, RandInt(0, 7, ticketId)),
		Fee1:     int64(RandInt(10000, 100000, ticketId)),
		Time2:    time.Now().AddDate(0, 0, RandInt(7, 14, ticketId)),
		Fee2:     int64(RandInt(20000, 150000, ticketId)),
		Time3:    time.Now().AddDate(0, 0, RandInt(14, 21, ticketId)),
		Fee3:     int64(RandInt(50000, 200000, ticketId)),
		Time4:    time.Now().AddDate(0, 0, RandInt(21, 28, ticketId)),
		Fee4:     int64(RandInt(80000, 250000, ticketId)),
		Time5:    time.Now().AddDate(0, 1, 0),
		Fee5:     int64(RandInt(100000, 250000, ticketId)),
	}

	rci1 := &commonModel.RefundAndChangeInfos{
		TicketId: ticketId,
		IsRefund: 1,
		Time1:    time.Now().AddDate(0, 0, RandInt(0, 3, ticketId)),
		Fee1:     int64(RandInt(10000, 100000, ticketId)),
		Time2:    time.Now().AddDate(0, 0, RandInt(3, 7, ticketId)),
		Fee2:     int64(RandInt(20000, 150000, ticketId)),
		Time3:    time.Now().AddDate(0, 0, RandInt(7, 10, ticketId)),
		Fee3:     int64(RandInt(50000, 200000, ticketId)),
		Time4:    time.Now().AddDate(0, 0, RandInt(10, 14, ticketId)),
		Fee4:     int64(RandInt(80000, 250000, ticketId)),
		Time5:    time.Now().AddDate(0, 1, 0),
		Fee5:     int64(RandInt(100000, 250000, ticketId)),
	}

	return []*commonModel.RefundAndChangeInfos{rci0, rci1}
}

func insertRCIs(c Config, minTicketID, maxTicketID int64) {
	var wg sync.WaitGroup
	rciModel := commonModel.NewRefundAndChangeInfosModel(sqlx.NewMysql(c.DB.DataSource), c.Cache)
	pool, _ := ants.NewPoolWithFunc(500, func(ticketId interface{}) {
		towTickets := genACoupleOfRandomRCI(ticketId.(int64))
		_, _ = rciModel.Insert(nil, towTickets[0])
		_, _ = rciModel.Insert(nil, towTickets[1])
		wg.Done()
	})
	defer pool.Release()

	for i := minTicketID; i <= maxTicketID; i++ {
		wg.Add(1)
		_ = pool.Invoke(i)
	}

	wg.Wait()
}
