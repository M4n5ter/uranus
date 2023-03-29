package logic

import (
	"context"
	"database/sql"
	"github.com/dtm-labs/dtm/client/dtmgrpc"
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"uranus/common/xerr"

	"uranus/app/stock/cmd/rpc/internal/svc"
	"uranus/app/stock/cmd/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type ReleaseStockByTicketIDLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewReleaseStockByTicketIDLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ReleaseStockByTicketIDLogic {
	return &ReleaseStockByTicketIDLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 通过 ticketID 释放锁定的库存
func (l *ReleaseStockByTicketIDLogic) ReleaseStockByTicketID(in *pb.ReleaseStockByTicketIDReq) (*pb.ReleaseStockResp, error) {

	if in.Num < 0 {
		return nil, status.Error(codes.Aborted, errors.Wrapf(ERRInvalidInput, "释放的库存数量不能为负数").Error())
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
		if space.LockedStock < 1 {
			return errors.Wrapf(xerr.NewErrMsg("可被释放的锁定库存不足"), "剩余锁定的库存量: %d", space.LockedStock)
		}
		space.LockedStock = space.LockedStock - in.Num
		err = l.svcCtx.SpacesModel.UpdateWithVersion(sqlx.NewSessionFromTx(tx), space)
		if err != nil {
			return errors.Wrapf(xerr.NewErrMsg("更新锁定库存失败"), "DBERR: %v", err)
		}
		return nil
	}); err != nil {
		return nil, status.Error(codes.Aborted, err.Error())
	}

	return &pb.ReleaseStockResp{}, nil
}
