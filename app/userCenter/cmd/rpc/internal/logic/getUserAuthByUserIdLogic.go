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

type GetUserAuthByUserIdLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetUserAuthByUserIdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserAuthByUserIdLogic {
	return &GetUserAuthByUserIdLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetUserAuthByUserIdLogic) GetUserAuthByUserId(in *pb.GetUserAuthByUserIdReq) (*pb.GetUserAuthByUserIdResp, error) {
	userAuth, err := l.svcCtx.UserAuthModel.FindOneByUserIdAuthType(in.UserId, in.AuthType)
	if err != nil && err != model.ErrNotFound {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DB_ERROR), "DB ERR when FindOneByUserIdAuthType, err: %v, authType: %s, userId: %d", err, in.AuthType, in.UserId)
	}
	if userAuth == nil {
		return nil, errors.Wrapf(ErrNotFoundUserAuthErr, "NotFoundUserAuth, authType: %s, userId: %d", in.AuthType, in.UserId)
	}
	return &pb.GetUserAuthByUserIdResp{UserAuth: &pb.UserAuth{
		Id:       userAuth.Id,
		UserId:   userAuth.UserId,
		AuthKey:  userAuth.AuthKey,
		AuthType: userAuth.AuthType,
	}}, nil
}
