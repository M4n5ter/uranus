package logic

import (
	"context"
	"github.com/pkg/errors"
	"uranus/app/userCenter/model"
	"uranus/common/xerr"

	"uranus/app/userCenter/cmd/rpc/internal/svc"
	"uranus/app/userCenter/cmd/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserMoneyLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetUserMoneyLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserMoneyLogic {
	return &GetUserMoneyLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetUserMoneyLogic) GetUserMoney(in *pb.GetUserMoneyReq) (*pb.GetUserMoneyResp, error) {
	if in.UserId < 1 {
		return nil, errors.Wrapf(xerr.NewErrMsg("非法输入"), "非法输入")
	}

	wallet, err := l.svcCtx.WalletModel.FindOneByUserId(in.UserId)
	if err != nil && err != model.ErrNotFound {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DB_ERROR), "DBERR in userCenter-rpc.GetUserMoney.l.svcCtx.WalletModel.FindOneByUserId, err: %v, userId: %d", err, in.UserId)
	}
	if wallet == nil {
		return nil, errors.Wrapf(xerr.NewErrMsg("用户没有钱包"), "Not Found wallet, userId: %d", in.UserId)
	}

	return &pb.GetUserMoneyResp{Money: wallet.Money}, nil
}
