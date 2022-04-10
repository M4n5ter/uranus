package logic

import (
	"context"
	"database/sql"
	"github.com/dtm-labs/dtm/dtmgrpc"
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"uranus/app/stock/cmd/rpc/internal/svc"
	"uranus/app/stock/cmd/rpc/pb"
	"uranus/common/xerr"

	"github.com/zeromicro/go-zero/core/logx"
)

type AddStockByTicketIDLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAddStockByTicketIDLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddStockByTicketIDLogic {
	return &AddStockByTicketIDLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 通过 ticketID 加库存
func (l *AddStockByTicketIDLogic) AddStockByTicketID(in *pb.AddStockByTicketIDReq) (*pb.AddStockResp, error) {
	if in.Num < 0 {
		return nil, errors.Wrapf(ERRInvalidInput, "加的库存数量不能为负数")
	}

	space, err := l.svcCtx.GetSpaceByTicketID(in.TicketID)
	if err != nil {
		return nil, err
	}

	barrier, err := dtmgrpc.BarrierFromGrpc(l.ctx)
	db, err := sqlx.NewMysql(l.svcCtx.Config.DB.DataSource).RawDB()
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	if err := barrier.CallWithDB(db, func(tx *sql.Tx) error {
		space.Surplus = space.Surplus + in.Num
		err = l.svcCtx.SpacesModel.UpdateWithVersion(sqlx.NewSessionFromTx(tx), space)
		if err != nil {
			return errors.Wrapf(xerr.NewErrMsg("更新库存失败"), "DBERR: %v", err)
		}
		return nil
	}); err != nil {
		return nil, status.Error(codes.Aborted, err.Error())
	}

	return &pb.AddStockResp{}, nil
}
