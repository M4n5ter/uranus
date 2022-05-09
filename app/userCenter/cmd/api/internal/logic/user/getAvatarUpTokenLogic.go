package user

import (
	"context"
	"github.com/qiniu/go-sdk/v7/auth/qbox"
	"github.com/qiniu/go-sdk/v7/storage"

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
	accessKey := l.svcCtx.Config.QiniuOSS.AccessKey
	secretKey := l.svcCtx.Config.QiniuOSS.SecretKey
	bucket := l.svcCtx.Config.QiniuOSS.Bucket
	mac := qbox.NewMac(accessKey, secretKey)
	putPolicy := storage.PutPolicy{
		Scope:      bucket,
		MimeLimit:  "image/jpeg;image/png",
		FsizeLimit: 10_485_760,
	}
	upToken := putPolicy.UploadToken(mac)

	return &types.GetAvatarUpTokenResp{UpToken: upToken}, nil
}
