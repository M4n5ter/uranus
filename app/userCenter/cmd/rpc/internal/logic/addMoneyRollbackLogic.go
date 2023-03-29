package logic

import (
	"context"
	"database/sql"
	"github.com/dtm-labs/dtm/client/dtmgrpc"
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"uranus/app/userCenter/model"
	"uranus/common/xerr"

	"uranus/app/userCenter/cmd/rpc/internal/svc"
	"uranus/app/userCenter/cmd/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type AddMoneyRollbackLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAddMoneyRollbackLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddMoneyRollbackLogic {
	return &AddMoneyRollbackLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *AddMoneyRollbackLogic) AddMoneyRollback(in *pb.AddMoneyReq) (*pb.AddMoneyResp, error) {

	// 检查用户是否有钱包
	wallet, err := l.svcCtx.WalletModel.FindOneByUserId(in.UserId)
	if err != nil && err != model.ErrNotFound {
		return nil, status.Error(codes.Internal, errors.Wrapf(xerr.NewErrCode(xerr.DB_ERROR), "DBERR: %v", err).Error())
	}

	if wallet == nil {
		return nil, status.Error(codes.Aborted, errors.Wrapf(xerr.NewErrMsg("该用户没有钱包"), "userID: %d", in.UserId).Error())
	}

	barrier, err := dtmgrpc.BarrierFromGrpc(l.ctx)
	db, err := sqlx.NewMysql(l.svcCtx.Config.DB.DataSource).RawDB()
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	if err := barrier.CallWithDB(db, func(tx *sql.Tx) error {
		// 回滚
		wallet.Money = wallet.Money - in.Money
		err = l.svcCtx.WalletModel.UpdateWithVersion(sqlx.NewSessionFromTx(tx), wallet)
		if err != nil {
			return errors.Wrapf(xerr.NewErrMsg("回滚加钱失败"), "DBERR: %v", err)
		}

		return nil
	}); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.AddMoneyResp{}, nil
}
