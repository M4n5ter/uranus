package svc

import (
	"github.com/zeromicro/go-zero/zrpc"
	"uranus/app/flightReservation/cmd/api/internal/config"
	"uranus/app/flightReservation/cmd/rpc/flightreservation"
)

type ServiceContext struct {
	Config                     config.Config
	FlightReservationRpcClient flightreservation.FlightReservation
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:                     c,
		FlightReservationRpcClient: flightreservation.NewFlightReservation(zrpc.MustNewClient(c.FlightReservation)),
	}
}
