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

type UploadAvatarLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUploadAvatarLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UploadAvatarLogic {
	return &UploadAvatarLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UploadAvatarLogic) UploadAvatar(req *types.UploadAvatarReq) (resp *types.UploadAvatarResp, err error) {
	userID := ctxdata.GetUidFromCtx(l.ctx)
	if userID < 1 {
		return nil, errors.Wrapf(xerr.NewErrMsg("用户不存在"), "用户不存在,ID:%d", userID)
	}

	_, err = l.svcCtx.UsercenterRpcClient.UploadAvatar(l.ctx, &usercenter.UploadAvatarReq{UserId: userID, Avatar: req.AvatarBasePath})
	if err != nil {
		return nil, errors.Wrapf(xerr.NewErrMsg("上传头像地址失败"), "RPC ERR : %+v", err)
	}

	return &types.UploadAvatarResp{}, nil
}
