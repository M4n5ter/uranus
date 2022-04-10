package svc

import (
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"uranus/app/stock/cmd/rpc/internal/config"
	"uranus/common/xerr"
	"uranus/commonModel"
)

type ServiceContext struct {
	Config       config.Config
	TicketsModel commonModel.TicketsModel
	SpacesModel  commonModel.SpacesModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:       c,
		TicketsModel: commonModel.NewTicketsModel(sqlx.NewMysql(c.DB.DataSource), c.Cache),
		SpacesModel:  commonModel.NewSpacesModel(sqlx.NewMysql(c.DB.DataSource), c.Cache),
	}
}

// GetSpaceByTicketID 通过 ticketID 得到 *spaces
func (s *ServiceContext) GetSpaceByTicketID(ticketID int64) (*commonModel.Spaces, error) {

	ticket, err := s.TicketsModel.FindOne(ticketID)
	if err != nil && err != commonModel.ErrNotFound {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DB_ERROR), "DBERR: %v", err)
	}
	if ticket == nil {
		return nil, errors.Wrapf(xerr.NewErrMsg("找不到对应票"), "ticketID: %d", ticketID)
	}

	space, err := s.SpacesModel.FindOne(ticket.SpaceId)
	if err != nil && err != commonModel.ErrNotFound {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DB_ERROR), "DBERR: %v", err)
	}
	if space == nil {
		return nil, errors.Wrapf(xerr.NewErrMsg("找不到对应舱位"), "spaceID: %d", ticket.SpaceId)
	}

	return space, nil
}
