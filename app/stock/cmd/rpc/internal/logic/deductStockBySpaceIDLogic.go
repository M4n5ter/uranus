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

type DeductStockBySpaceIDLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeductStockBySpaceIDLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeductStockBySpaceIDLogic {
	return &DeductStockBySpaceIDLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 通过 spaceID 扣库存
func (l *DeductStockBySpaceIDLogic) DeductStockBySpaceID(in *pb.DeductStockBySpaceIDReq) (*pb.DeductStockResp, error) {

	if in.Num < 0 {
		return nil, status.Error(codes.Aborted, errors.Wrapf(ERRInvalidInput, "扣的库存数量不能为负数").Error())
	}

	space, err := l.svcCtx.SpacesModel.FindOne(in.SpaceID)
	if err != nil && err != commonModel.ErrNotFound {
		return nil, status.Error(codes.Aborted, errors.Wrapf(xerr.NewErrCode(xerr.DB_ERROR), "DBERR: %v", err).Error())
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
		if space.Surplus-space.LockedStock < 0 {
			return errors.Wrapf(xerr.NewErrMsg("库存不足"), "")
		}
		space.Surplus, space.LockedStock = space.Surplus-in.Num, space.LockedStock-in.Num
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
