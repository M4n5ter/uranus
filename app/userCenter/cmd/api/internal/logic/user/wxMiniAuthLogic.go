package user

import (
	"context"
	"github.com/jinzhu/copier"
	"github.com/pkg/errors"
	"github.com/silenceper/wechat/v2/cache"
	miniConfig "github.com/silenceper/wechat/v2/miniprogram/config"
	"uranus/app/uranusAuth/cmd/rpc/auth"
	"uranus/app/userCenter/cmd/rpc/userCenter"
	"uranus/app/userCenter/model"
	"uranus/common/tool"
	"uranus/common/xerr"

	"github.com/silenceper/wechat/v2"
	"github.com/zeromicro/go-zero/core/logx"
	"uranus/app/userCenter/cmd/api/internal/svc"
	"uranus/app/userCenter/cmd/api/internal/types"
)

var ErrMiniAuthFailedErr = xerr.NewErrMsg("授权失败")

type WxMiniAuthLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewWxMiniAuthLogic(ctx context.Context, svcCtx *svc.ServiceContext) WxMiniAuthLogic {
	return WxMiniAuthLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *WxMiniAuthLogic) WxMiniAuth(req *types.WXMiniAuthReq) (resp *types.WXMiniAuthResp, err error) {
	// 微信授权
	miniProgram := wechat.NewWechat().GetMiniProgram(&miniConfig.Config{
		AppID:     l.svcCtx.Config.WxMiniConf.AppId,
		AppSecret: l.svcCtx.Config.WxMiniConf.Secret,
		Cache:     cache.NewMemory(),
	})
	authResult, err := miniProgram.GetAuth().Code2Session(req.Code)
	if err != nil || authResult.ErrCode != 0 || authResult.OpenID == "" {
		return nil, errors.Wrapf(ErrMiniAuthFailedErr, "发起授权请求失败,code: %s, err: %v, authResult: %+v", req.Code, err, authResult)
	}
	// 解析小程序返回的数据
	userData, err := miniProgram.GetEncryptor().Decrypt(authResult.SessionKey, req.EncryptedData, req.IV)
	if err != nil {
		return nil, errors.Wrapf(ErrMiniAuthFailedErr, "解析数据失败, err: %v, authResult: %+v, req: %+v", err, authResult, req)
	}
	// 绑定用户(未绑定)或者登录(已绑定)
	var userId int64
	rpcResp, err := l.svcCtx.UsercenterRpcClient.GetUserAuthByAuthKey(l.ctx, &userCenter.GetUserAuthByAuthKeyReq{
		AuthKey:  authResult.OpenID,
		AuthType: model.UserAuthTypeSmallWX,
	})
	if err != nil {
		return nil, errors.Wrapf(ErrMiniAuthFailedErr, "UsercenterRpc.GetUserAuthByAuthKey err: %v, authResult : %+v", err, authResult)
	}
	if rpcResp.UserAuth == nil || rpcResp.UserAuth.Id == 0 {
		// 用户未绑定，需要进行绑定
		mobile := userData.PhoneNumber
		// 昵称为 4 位随机字符串加上手机号后四位
		nickname := tool.Krand(4, tool.KC_RAND_KIND_ALL) + mobile[8:]
		registerResp, err := l.svcCtx.UsercenterRpcClient.Register(l.ctx, &userCenter.RegisterReq{
			Mobile:   mobile,
			Nickname: nickname,
			Sex:      model.Unknown,
			AuthKey:  authResult.OpenID,
			AuthType: model.UserAuthTypeSmallWX,
		})
		if err != nil {
			return nil, errors.Wrapf(ErrMiniAuthFailedErr, "UsercenterRpc.Register err :%v, authResult : %+v", err, authResult)
		}
		_ = copier.Copy(resp, registerResp)
	} else {
		// 用户已经绑定，登录:直接授权返回token
		userId = rpcResp.UserAuth.UserId
		tokenResp, err := l.svcCtx.AuthRpcClient.GenerateToken(l.ctx, &auth.GenerateTokenReq{UserId: userId})
		if err != nil {
			return nil, errors.Wrapf(ErrMiniAuthFailedErr, "authRpc.GenerateTokenReq err: %v, userId: %d", err, userId)
		}
		_ = copier.Copy(resp, tokenResp)
	}
	return
}
