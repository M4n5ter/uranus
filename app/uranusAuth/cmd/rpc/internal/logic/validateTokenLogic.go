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

var ErrValidateErr = xerr.NewErrMsg("验证 token 失败")

type ValidateTokenLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewValidateTokenLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ValidateTokenLogic {
	return &ValidateTokenLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// ValidateToken 只对 usercenter 和 uranusAuth api 开放
func (l *ValidateTokenLogic) ValidateToken(in *pb.ValidateTokenReq) (*pb.ValidateTokenResp, error) {
	userTokenKey := fmt.Sprintf(globalkey.CacheUserTokenKey, in.UserId)
	token, err := l.svcCtx.RedisClient.Get(userTokenKey)
	if err != nil {
		return nil, errors.Wrapf(ErrValidateErr, "get CacheUserTokenKey failed, err: %v", err)
	}
	if token != in.Token {
		return nil, errors.Wrapf(ErrValidateErr, "token is invalid, CacheUserToken: %s, invalidToken: %s, userId: %d", token, in.Token, in.UserId)
	}

	return &pb.ValidateTokenResp{Ok: true}, nil
}
