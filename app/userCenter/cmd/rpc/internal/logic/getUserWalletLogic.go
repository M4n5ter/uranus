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

type GetUserWalletLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetUserWalletLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserWalletLogic {
	return &GetUserWalletLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetUserWalletLogic) GetUserWallet(in *pb.GetUserWalletReq) (*pb.GetUserWalletResp, error) {
	if in.UserId < 1 {
		return nil, errors.Wrapf(xerr.NewErrMsg("非法用户"), "")
	}

	wallet, err := l.svcCtx.WalletModel.FindOneByUserId(in.UserId)
	if err != nil && err != model.ErrNotFound {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DB_ERROR), "DBERR: %+v", err)
	}

	if wallet == nil {
		return nil, errors.Wrapf(xerr.NewErrMsg("用户没有钱包"), "")
	}

	return &pb.GetUserWalletResp{Balance: wallet.Money}, nil
}
