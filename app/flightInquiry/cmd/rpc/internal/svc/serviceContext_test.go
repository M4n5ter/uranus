package svc

import (
	"flag"
	"testing"
)

var configFile = flag.String("f", "../../etc/flightInquiry.yaml", "the config file")

func TestCombineAllInfos(t *testing.T) {
	//flag.Parse()
	//
	//var c config.Config
	//conf.MustLoad(*configFile, &c)
	//fli := &commonModel.FlightInfos{
	//	Id:           1,
	//	DelState:     0,
	//	Version:      0,
	//	FlightNumber: "KY8231",
	//}
	//ret, err := NewServiceContext(c).CombineAllInfos([]*commonModel.FlightInfos{fli})
	//if err != nil {
	//	t.Error(err)
	//	t.FailNow()
	//}
	//if len(ret) != 1 {
	//	t.Logf("实际条数: %d", len(ret))
	//	t.Fail()
	//}

}
