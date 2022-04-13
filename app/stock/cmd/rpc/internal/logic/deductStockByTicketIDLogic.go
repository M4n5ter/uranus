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

type DeductStockByTicketIDLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeductStockByTicketIDLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeductStockByTicketIDLogic {
	return &DeductStockByTicketIDLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 通过 ticketID 扣库存
func (l *DeductStockByTicketIDLogic) DeductStockByTicketID(in *pb.DeductStockByTicketIDReq) (*pb.DeductStockResp, error) {

	if in.Num < 0 {
		return nil, status.Error(codes.Aborted, errors.Wrapf(ERRInvalidInput, "扣的库存数量不能为负数").Error())
	}

	space, err := l.svcCtx.GetSpaceByTicketID(in.TicketID)
	if err != nil {
		if err.(*xerr.CodeError).GetErrCode() == xerr.DB_ERROR {
			return nil, status.Error(codes.Internal, err.Error())
		}
		return nil, status.Error(codes.Aborted, err.Error())
	}

	barrier, err := dtmgrpc.BarrierFromGrpc(l.ctx)
	db, err := sqlx.NewMysql(l.svcCtx.Config.DB.DataSource).RawDB()
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	if err := barrier.CallWithDB(db, func(tx *sql.Tx) error {
		if space.Surplus-space.LockedStock < 0 {
			return errors.Wrapf(xerr.NewErrMsg("库存不足"), "")
		}
		space.Surplus = space.Surplus - in.Num
		err = l.svcCtx.SpacesModel.UpdateWithVersion(sqlx.NewSessionFromTx(tx), space)
		if err != nil {
			return errors.Wrapf(xerr.NewErrMsg("更新库存失败"), "DBERR: %v", err)
		}
		return nil
	}); err != nil {
		return nil, status.Error(codes.Aborted, err.Error())
	}

	return &pb.DeductStockResp{}, nil
}
