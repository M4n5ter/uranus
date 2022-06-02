package user

import (
	"context"
	"uranus/app/userCenter/cmd/rpc/usercenter"
	"uranus/common/ctxdata"

	"uranus/app/userCenter/cmd/api/internal/svc"
	"uranus/app/userCenter/cmd/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserWalletLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetUserWalletLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserWalletLogic {
	return &GetUserWalletLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetUserWalletLogic) GetUserWallet(req *types.GetUserWalletReq) (resp *types.GetUserWalletResp, err error) {
	userID := ctxdata.GetUidFromCtx(l.ctx)
	rpcResp, err := l.svcCtx.UsercenterRpcClient.GetUserWallet(l.ctx, &usercenter.GetUserWalletReq{UserId: userID})
	if err != nil {
		return nil, err
	}

	resp = &types.GetUserWalletResp{Balance: rpcResp.Balance}
	return
}
