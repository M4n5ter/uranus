package user

import (
	"context"
	"github.com/jinzhu/copier"
	"github.com/pkg/errors"
	"uranus/app/userCenter/model"
	"uranus/app/usercenter/cmd/api/internal/svc"
	"uranus/app/usercenter/cmd/api/internal/types"
	"uranus/app/usercenter/cmd/rpc/pb"
	"uranus/common/xerr"

	"github.com/zeromicro/go-zero/core/logx"
)

var ErrEmptyMobileErr = xerr.NewErrMsg("手机号不能为空")
var ErrRegisterErr = xerr.NewErrMsg("注册失败")

type RegisterLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) RegisterLogic {
	return RegisterLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RegisterLogic) Register(req *types.RegisterReq) (resp *types.RegisterResp, err error) {
	if len(req.Mobile) == 0 {
		return nil, errors.Wrapf(ErrEmptyMobileErr, "emptyMobile err, req: %+v", req)
	}
	rpcResp, err := l.svcCtx.UsercenterRpcClient.Register(l.ctx, &pb.RegisterReq{
		Mobile:   req.Mobile,
		Password: req.Password,
		AuthKey:  req.Mobile,
		AuthType: model.UserAuthTypeSystem,
		Sex:      req.Sex,
		Nickname: req.Nickname,
	})
	if err != nil {
		return nil, errors.Wrapf(ErrRegisterErr, "req: %+v", req)
	}
	_ = copier.Copy(resp, rpcResp)
	return
}
