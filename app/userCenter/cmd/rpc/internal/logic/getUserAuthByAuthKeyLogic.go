package logic

import (
	"context"
	"github.com/pkg/errors"
	"uranus/app/userCenter/model"
	"uranus/common/xerr"

	"uranus/app/usercenter/cmd/rpc/internal/svc"
	"uranus/app/usercenter/cmd/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

var ErrNotFoundUserAuthErr = xerr.NewErrMsg("未找到 UserAuth")

type GetUserAuthByAuthKeyLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetUserAuthByAuthKeyLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserAuthByAuthKeyLogic {
	return &GetUserAuthByAuthKeyLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetUserAuthByAuthKeyLogic) GetUserAuthByAuthKey(in *pb.GetUserAuthByAuthKeyReq) (*pb.GetUserAuthByAuthKeyResp, error) {
	userAuth, err := l.svcCtx.UserAuthModel.FindOneByAuthTypeAuthKey(in.AuthType, in.AuthKey)
	if err != nil && err != model.ErrNotFound {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DB_ERROR), "DB ERR when FindOneByAuthTypeAuthKey, err: %v, authType: %s, authKey: %s", err, in.AuthType, in.AuthKey)
	}
	if userAuth == nil {
		return nil, errors.Wrapf(ErrNotFoundUserAuthErr, "NotFoundUserAuth, authType: %s, authKey: %s", in.AuthType, in.AuthKey)
	}
	return &pb.GetUserAuthByAuthKeyResp{UserAuth: &pb.UserAuth{
		Id:       userAuth.Id,
		UserId:   userAuth.UserId,
		AuthKey:  userAuth.AuthKey,
		AuthType: userAuth.AuthType,
	}}, nil
}
