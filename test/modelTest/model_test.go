package modelTest

import (
	"fmt"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"testing"
	"time"
	"uranus/commonModel"
)

type Config struct {
	DB struct {
		DataSource string
	}
	Cache cache.CacheConf
}

func TestFindTransferFlightsByPlace(t *testing.T) {
	var c Config
	conf.MustLoad("D:\\Project\\uranus\\test\\modelTest\\etc\\genData.yaml", &c)
	flightInfosModel := commonModel.NewFlightInfosModel(sqlx.NewMysql(c.DB.DataSource), c.Cache)
	transfers, err := flightInfosModel.FindTransferFlightsByPlace(flightInfosModel.RowBuilder(), "白云国际机场T2", "萧山国际机场T1", time.Date(2022, 3, 22, 0, 0, 0, 0, time.UTC))
	if err != nil {
		t.Error(err)
	}
	fmt.Printf("%+v", transfers)
}
