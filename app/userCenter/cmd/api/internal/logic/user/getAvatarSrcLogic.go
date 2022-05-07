package user

import (
	"context"
	"github.com/pkg/errors"
	"uranus/app/userCenter/cmd/rpc/usercenter"
	"uranus/common/ctxdata"
	"uranus/common/xerr"

	"uranus/app/userCenter/cmd/api/internal/svc"
	"uranus/app/userCenter/cmd/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetAvatarSrcLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetAvatarSrcLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetAvatarSrcLogic {
	return &GetAvatarSrcLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetAvatarSrcLogic) GetAvatarSrc(req *types.GetAvatarSrcReq) (resp *types.GetAvatarSrcResp, err error) {
	userID := ctxdata.GetUidFromCtx(l.ctx)
	if userID < 1 {
		return nil, errors.Wrapf(xerr.NewErrMsg("用户不存在"), "用户不存在,ID:%d", userID)
	}

	getAvatarSrcResp, err := l.svcCtx.UsercenterRpcClient.GetAvatarSrc(l.ctx, &usercenter.GetAvatarSrcReq{UserId: userID})
	if err != nil {
		return nil, errors.Wrapf(xerr.NewErrMsg("获取头像地址失败"), "RPC ERR : %+v", err)
	}

	return &types.GetAvatarSrcResp{AvatarSrc: getAvatarSrcResp.Avatar}, nil
}
