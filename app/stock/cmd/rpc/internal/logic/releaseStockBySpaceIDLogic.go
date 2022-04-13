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

type ReleaseStockBySpaceIDLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewReleaseStockBySpaceIDLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ReleaseStockBySpaceIDLogic {
	return &ReleaseStockBySpaceIDLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 通过 spaceID 释放锁定的库存
func (l *ReleaseStockBySpaceIDLogic) ReleaseStockBySpaceID(in *pb.ReleaseStockBySpaceIDReq) (*pb.ReleaseStockResp, error) {

	if in.Num < 0 {
		return nil, status.Error(codes.Aborted, errors.Wrapf(ERRInvalidInput, "释放的库存数量不能为负数").Error())
	}

	space, err := l.svcCtx.SpacesModel.FindOne(in.SpaceID)
	if err != nil && err != commonModel.ErrNotFound {
		return nil, status.Error(codes.Internal, errors.Wrapf(xerr.NewErrCode(xerr.DB_ERROR), "DBERR: %v", err).Error())
	}
	if space == nil {
		return nil, status.Error(codes.Aborted, errors.Wrapf(xerr.NewErrMsg("找不到对应舱位"), "spaceID: %d", in.SpaceID).Error())
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
