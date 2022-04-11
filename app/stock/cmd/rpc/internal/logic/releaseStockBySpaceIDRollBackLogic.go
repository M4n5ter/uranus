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
	"uranus/commonModel"

	"uranus/app/stock/cmd/rpc/internal/svc"
	"uranus/app/stock/cmd/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type ReleaseStockBySpaceIDRollBackLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewReleaseStockBySpaceIDRollBackLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ReleaseStockBySpaceIDRollBackLogic {
	return &ReleaseStockBySpaceIDRollBackLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 通过 spaceID 释放锁定的库存 rollback
func (l *ReleaseStockBySpaceIDRollBackLogic) ReleaseStockBySpaceIDRollBack(in *pb.ReleaseStockBySpaceIDReq) (*pb.ReleaseStockResp, error) {

	space, err := l.svcCtx.SpacesModel.FindOne(in.SpaceID)
	if err != nil && err != commonModel.ErrNotFound {
		return nil, status.Error(codes.Internal, err.Error())
	}
	if space == nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	barrier, err := dtmgrpc.BarrierFromGrpc(l.ctx)
	db, err := sqlx.NewMysql(l.svcCtx.Config.DB.DataSource).RawDB()
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	if err := barrier.CallWithDB(db, func(tx *sql.Tx) error {
		space.LockedStock = space.LockedStock + in.Num
		err = l.svcCtx.SpacesModel.UpdateWithVersion(sqlx.NewSessionFromTx(tx), space)
		if err != nil {
			return errors.Wrapf(xerr.NewErrMsg("回滚锁定库存失败"), "DBERR: %v", err)
		}
		return nil
	}); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.ReleaseStockResp{}, nil
}
