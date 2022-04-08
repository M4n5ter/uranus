package user

import (
	"context"
	"github.com/jinzhu/copier"
	"uranus/app/userCenter/cmd/rpc/userCenter"
	"uranus/common/ctxdata"

	"uranus/app/userCenter/cmd/api/internal/svc"
	"uranus/app/userCenter/cmd/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DetailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) DetailLogic {
	return DetailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// Detail 获取当前登录用户的信息
func (l *DetailLogic) Detail(req *types.UserInfoReq) (resp *types.UserInfoResp, err error) {
	userId := ctxdata.GetUidFromCtx(l.ctx)
	userInfoResp, err := l.svcCtx.UsercenterRpcClient.GetUserInfo(l.ctx, &userCenter.GetUserInfoReq{Id: userId})
	if err != nil {
		return nil, err
	}
	userInfo := userInfoResp.User
	resp = &types.UserInfoResp{}
	_ = copier.Copy(resp, userInfo)
	return
}
