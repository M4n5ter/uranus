package flightInquiry

import (
	"context"

	"uranus/app/flightInquiry/cmd/api/internal/svc"
	"uranus/app/flightInquiry/cmd/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetFlightsByFlightNumberLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetFlightsByFlightNumberLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetFlightsByFlightNumberLogic {
	return &GetFlightsByFlightNumberLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetFlightsByFlightNumberLogic) GetFlightsByFlightNumber(req *types.GetFlightsByFlightNumberReq) (resp *types.GetFlightsByFlightNumberResp, err error) {
	// todo: add your logic here and delete this line

	return
}
