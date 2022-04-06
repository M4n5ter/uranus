package user

import (
	"context"
	"github.com/pkg/errors"
	"uranus/app/userCenter/cmd/rpc/userCenter"
	"uranus/app/userCenter/model"
	"uranus/common/xerr"

	"uranus/app/userCenter/cmd/api/internal/svc"
	"uranus/app/userCenter/cmd/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

var ErrPlatformErr = xerr.NewErrMsg("平台类型错误，该接口不支持微信小程序")
var ErrInvalidMobileOrPassword = xerr.NewErrMsg("账号或密码无效")

type LoginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) LoginLogic {
	return LoginLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// Login 是本地系统平台用户登录渠道，微信小程序登录渠道是 wxMiniAuth
func (l *LoginLogic) Login(req *types.LoginReq) (resp *types.LoginResp, err error) {
	if req.AuthType == model.UserAuthTypeSmallWX {
		return nil, errors.Wrapf(ErrPlatformErr, "err platform: %v, req: %+v", req.AuthType, req)
	}
	if len(req.Mobile) == 0 {
		return nil, errors.Wrapf(ErrInvalidMobileOrPassword, "empty mobile err, req: %+v", req)
	}
	if len(req.Password) == 0 {
		return nil, errors.Wrapf(ErrInvalidMobileOrPassword, "empty password err, req: %+v", req)
	}
	loginResp, err := l.svcCtx.UsercenterRpcClient.Login(l.ctx, &userCenter.LoginReq{
		AuthKey:  req.Mobile,
		AuthType: req.AuthType,
		Password: req.Password,
	})
	if err != nil {
		return nil, err
	}
	resp = &types.LoginResp{}
	resp.AccessToken = loginResp.AccessToken
	resp.AccessExpire = loginResp.AccessExpire
	resp.RefreshAfter = loginResp.RefreshAfter
	return
}
