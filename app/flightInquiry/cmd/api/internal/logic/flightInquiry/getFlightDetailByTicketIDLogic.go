package flightInquiry

import (
	"context"

	"uranus/app/flightInquiry/cmd/api/internal/svc"
	"uranus/app/flightInquiry/cmd/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetFlightDetailByTicketIDLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetFlightDetailByTicketIDLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetFlightDetailByTicketIDLogic {
	return &GetFlightDetailByTicketIDLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetFlightDetailByTicketIDLogic) GetFlightDetailByTicketID(req *types.GetFlightDetailByTicketIDReq) (resp *types.GetFlightDetailByTicketIDResp, err error) {
	// todo: add your logic here and delete this line

	return
}
