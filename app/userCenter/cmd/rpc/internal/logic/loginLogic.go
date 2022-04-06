package logic

import (
	"context"
	"github.com/jinzhu/copier"
	"github.com/pkg/errors"
	"uranus/app/uranusAuth/cmd/rpc/auth"
	"uranus/app/userCenter/model"
	"uranus/common/xerr"

	"uranus/app/userCenter/cmd/rpc/internal/svc"
	"uranus/app/userCenter/cmd/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

var ErrLoginErr = xerr.NewErrMsg("登录失败，请确认账号和密码是否正确")

type LoginLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *LoginLogic) Login(in *pb.LoginReq) (*pb.LoginResp, error) {
	user, err := l.svcCtx.UserModel.FindOneByMobile(in.AuthKey)
	if err != nil && err != model.ErrNotFound {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DB_ERROR), "DB ERR When finding one user by mobile, mobile: %s, err: %v", in.AuthType, err)
	}
	if user == nil {
		return nil, errors.Wrapf(ErrLoginErr, "mobile doesn't exist, req: %+v", in)
	}
	resp := &pb.LoginResp{}
	if user.Password == in.Password {
		authResp, err := l.svcCtx.AuthRpcClient.GenerateToken(l.ctx, &auth.GenerateTokenReq{UserId: user.Id})
		if err != nil {
			return nil, err
		}
		_ = copier.Copy(resp, authResp)
	}

	return resp, nil
}
