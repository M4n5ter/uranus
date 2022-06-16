package user

import (
	"context"
	"github.com/pkg/errors"
	"github.com/qiniu/go-sdk/v7/auth/qbox"
	"github.com/qiniu/go-sdk/v7/storage"
	"github.com/zeromicro/go-zero/core/hash"
	"uranus/app/userCenter/cmd/rpc/usercenter"
	"uranus/common/ctxdata"
	"uranus/common/xerr"

	"uranus/app/userCenter/cmd/api/internal/svc"
	"uranus/app/userCenter/cmd/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetAvatarUpTokenLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetAvatarUpTokenLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetAvatarUpTokenLogic {
	return &GetAvatarUpTokenLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetAvatarUpTokenLogic) GetAvatarUpToken(req *types.GetAvatarUpTokenReq) (resp *types.GetAvatarUpTokenResp, err error) {
	userid := ctxdata.GetUidFromCtx(l.ctx)
	if userid == 0 {
		return nil, errors.Wrapf(xerr.NewErrMsg("非法用户"), "用户 ID 为 0！")
	}

	userinfoResp, err := l.svcCtx.UsercenterRpcClient.GetUserInfo(l.ctx, &usercenter.GetUserInfoReq{Id: userid})
	if err != nil {
		return nil, err
	}

	nickname, mobile := userinfoResp.User.Nickname, userinfoResp.User.Mobile
	key := hash.Md5Hex([]byte(nickname + "_" + mobile))

	accessKey := l.svcCtx.Config.QiniuOSS.AccessKey
	secretKey := l.svcCtx.Config.QiniuOSS.SecretKey
	bucket := l.svcCtx.Config.QiniuOSS.Bucket
	mac := qbox.NewMac(accessKey, secretKey)
	putPolicy := storage.PutPolicy{
		Scope:        bucket,
		MimeLimit:    "image/jpeg;image/png",
		FsizeLimit:   10_485_760,
		ForceSaveKey: true,
		SaveKey:      key,
	}
	upToken := putPolicy.UploadToken(mac)

	return &types.GetAvatarUpTokenResp{UpToken: upToken}, nil
}
