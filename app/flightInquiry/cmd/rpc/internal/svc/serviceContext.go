package svc

import (
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"uranus/app/flightInquiry/cmd/rpc/internal/config"
	"uranus/model"
)

type ServiceContext struct {
	Config                    config.Config
	Flights                   model.FlightsModel
	FlightInfosModel          model.FlightInfosModel
	SpacesModel               model.SpacesModel
	TicketsModel              model.TicketsModel
	RefundAndChangeInfosModel model.RefundAndChangeInfosModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:                    c,
		Flights:                   model.NewFlightsModel(sqlx.NewMysql(c.DB.DataSource), c.Cache),
		FlightInfosModel:          model.NewFlightInfosModel(sqlx.NewMysql(c.DB.DataSource), c.Cache),
		SpacesModel:               model.NewSpacesModel(sqlx.NewMysql(c.DB.DataSource), c.Cache),
		TicketsModel:              model.NewTicketsModel(sqlx.NewMysql(c.DB.DataSource), c.Cache),
		RefundAndChangeInfosModel: model.NewRefundAndChangeInfosModel(sqlx.NewMysql(c.DB.DataSource), c.Cache),
	}
}
