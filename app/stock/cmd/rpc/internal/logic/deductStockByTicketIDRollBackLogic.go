package logic

import (
	"context"
	"database/sql"
	"github.com/dtm-labs/dtm/dtmgrpc"
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"uranus/common/xerr"

	"uranus/app/stock/cmd/rpc/internal/svc"
	"uranus/app/stock/cmd/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeductStockByTicketIDRollBackLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeductStockByTicketIDRollBackLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeductStockByTicketIDRollBackLogic {
	return &DeductStockByTicketIDRollBackLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 通过 ticketID 扣库存 rollback
func (l *DeductStockByTicketIDRollBackLogic) DeductStockByTicketIDRollBack(in *pb.DeductStockByTicketIDReq) (*pb.DeductStockResp, error) {

	space, err := l.svcCtx.GetSpaceByTicketID(in.TicketID)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	barrier, err := dtmgrpc.BarrierFromGrpc(l.ctx)
	db, err := sqlx.NewMysql(l.svcCtx.Config.DB.DataSource).RawDB()
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	if err := barrier.CallWithDB(db, func(tx *sql.Tx) error {
		space.LockedStock = space.LockedStock - in.Num
		space.Surplus = space.Surplus + in.Num
		err = l.svcCtx.SpacesModel.UpdateWithVersion(sqlx.NewSessionFromTx(tx), space)
		if err != nil {
			return errors.Wrapf(xerr.NewErrMsg("更新库存失败"), "DBERR: %v", err)
		}
		return nil
	}); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.DeductStockResp{}, nil
}
