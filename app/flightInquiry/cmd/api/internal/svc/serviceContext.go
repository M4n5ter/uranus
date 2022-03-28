package svc

import (
	"github.com/zeromicro/go-zero/zrpc"
	"uranus/app/flightInquiry/cmd/api/internal/config"
	"uranus/app/flightInquiry/cmd/rpc/flightinquiry"
)

type ServiceContext struct {
	Config              config.Config
	FlightInquiryClient flightinquiry.FlightInquiry
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:              c,
		FlightInquiryClient: flightinquiry.NewFlightInquiry(zrpc.MustNewClient(c.FlightInquiryConf)),
	}
}
