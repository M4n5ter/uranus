package logic

import (
	"context"
	"github.com/jinzhu/copier"
	"github.com/pkg/errors"
	"uranus/app/usercenter/model"
	"uranus/common/xerr"

	"uranus/app/usercenter/cmd/rpc/internal/svc"
	"uranus/app/usercenter/cmd/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

var ErrUnknownErr = xerr.NewErrMsg("未知用户")

type GetUserInfoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetUserInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserInfoLogic {
	return &GetUserInfoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetUserInfoLogic) GetUserInfo(in *pb.GetUserInfoReq) (*pb.GetUserInfoResp, error) {
	userInfoResp, err := l.svcCtx.UserModel.FindOne(in.Id)
	if err != nil && err != model.ErrNotFound {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DB_ERROR), "DB ERR when finding one in user table, id: %d, err: %v", in.Id, err)
	}
	if userInfoResp == nil {
		return nil, errors.Wrapf(ErrUnknownErr, "Err: unknown user but it has an userId, id: %d", in.Id)
	}
	resp := &pb.GetUserInfoResp{}
	_ = copier.Copy(resp, userInfoResp)
	return resp, nil
}
