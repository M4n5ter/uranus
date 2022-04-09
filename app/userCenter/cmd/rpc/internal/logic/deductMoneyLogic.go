package logic

import (
	"context"
	"database/sql"
	"github.com/dtm-labs/dtm/dtmgrpc"
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

type DeductMoneyLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeductMoneyLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeductMoneyLogic {
	return &DeductMoneyLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DeductMoneyLogic) DeductMoney(in *pb.DeductMoneyReq) (*pb.DeductMoneyResp, error) {
	// 检查输入
	if in.Money < 0 {
		return nil, errors.Wrapf(xerr.NewErrMsg("扣的钱不能小于 0 ，加钱请使用隔壁 rpc"), "")
	}

	// 检查用户是否存在
	_, err := l.svcCtx.UserModel.FindOne(in.UserId)
	if err != nil {
		if err == model.ErrNotFound {
			return nil, errors.Wrapf(xerr.NewErrMsg("用户不存在"), "Not Found userID: %d", in.UserId)
		}
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DB_ERROR), "DBERR: %v", err)
	}

	// 检查用户是否有钱包
	wallet, err := l.svcCtx.WalletModel.FindOneByUserId(in.UserId)
	if err != nil && err != model.ErrNotFound {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DB_ERROR), "DBERR: %v", err)
	}

	if wallet == nil {
		return nil, errors.Wrapf(xerr.NewErrMsg("该用户没有钱包"), "userID: %d", in.UserId)
	}

	barrier, err := dtmgrpc.BarrierFromGrpc(l.ctx)
	db, err := sqlx.NewMysql(l.svcCtx.Config.DB.DataSource).RawDB()
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	if err := barrier.CallWithDB(db, func(tx *sql.Tx) error {
		// 扣钱(可以为负，如要判断是否有足够余额，请判断后再调用 扣钱 rpc)
		wallet.Money = wallet.Money - in.Money
		err = l.svcCtx.WalletModel.UpdateWithVersion(sqlx.NewSessionFromTx(tx), wallet)
		if err != nil {
			return errors.Wrapf(xerr.NewErrMsg("扣钱失败"), "DBERR: %v", err)
		}
		return nil
	}); err != nil {
		return nil, status.Error(codes.Aborted, err.Error())
	}
	return &pb.DeductMoneyResp{}, nil
}
