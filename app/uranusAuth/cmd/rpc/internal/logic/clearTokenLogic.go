package logic

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"uranus/common/globalkey"
	"uranus/common/xerr"

	"uranus/app/uranusAuth/cmd/rpc/internal/svc"
	"uranus/app/uranusAuth/cmd/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

var ErrClearTokenErr = xerr.NewErrMsg("清除 token 失败")

type ClearTokenLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewClearTokenLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ClearTokenLogic {
	return &ClearTokenLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// ClearToken 清除 Token 服务只对 userCenter 开放
func (l *ClearTokenLogic) ClearToken(in *pb.ClearTokenReq) (*pb.ClearTokenResp, error) {
	userTokenKey := fmt.Sprintf(globalkey.CacheUserTokenKey, in.UserId)
	if _, err := l.svcCtx.RedisClient.Del(userTokenKey); err != nil {
		return nil, errors.Wrapf(ErrClearTokenErr, "delete userTokenKey cache failed: userId: %v, err: %v", in.UserId, err)
	}
	return &pb.ClearTokenResp{
		Ok: true,
	}, nil
}
