package logic

import (
	"context"
	"github.com/pkg/errors"
	"github.com/qiniu/go-sdk/v7/storage"
	"time"
	userCenterModel "uranus/app/userCenter/model"
	"uranus/common/xerr"

	"uranus/app/userCenter/cmd/rpc/internal/svc"
	"uranus/app/userCenter/cmd/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetAvatarSrcLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetAvatarSrcLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetAvatarSrcLogic {
	return &GetAvatarSrcLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetAvatarSrcLogic) GetAvatarSrc(in *pb.GetAvatarSrcReq) (*pb.GetAvatarSrcResp, error) {
	user, err := l.svcCtx.UserModel.FindOne(in.UserId)
	if err != nil && err != userCenterModel.ErrNotFound {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DB_ERROR), "DBERR: %+v", err)
	}

	if user == nil {
		return nil, errors.Wrapf(xerr.NewErrMsg("用户不存在"), "用户不存在,userID: %d", in.UserId)
	}

	if len(user.Avatar) == 0 {
		//没有头像
		return nil, errors.Wrapf(xerr.NewErrMsg("用户没有头像"), "")
	}

	mac := l.svcCtx.GenQiniuOSSMAC()
	//生成 1 小时有效的链接地址
	domain := l.svcCtx.Config.QiniuOSS.Domain
	key := user.Avatar
	deadline := time.Now().Add(time.Second * 3600).Unix() //1小时有效期
	privateAccessURL := storage.MakePrivateURL(mac, domain, key, deadline)

	return &pb.GetAvatarSrcResp{Avatar: privateAccessURL}, nil
}
