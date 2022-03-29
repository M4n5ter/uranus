package logic

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
	"time"
	"uranus/app/uranusAuth/cmd/rpc/internal/svc"
	"uranus/app/uranusAuth/cmd/rpc/pb"
	"uranus/common/globalkey"
	"uranus/common/jwt"
	"uranus/common/xerr"
)

var ErrGenerateTokenError = xerr.NewErrMsg("生成token失败")

type GenerateTokenLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGenerateTokenLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GenerateTokenLogic {
	return &GenerateTokenLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// GenerateToken 生成 Token 服务只对 userCenter 开放
func (l *GenerateTokenLogic) GenerateToken(in *pb.GenerateTokenReq) (*pb.GenerateTokenResp, error) {
	now := time.Now().Unix()
	accessExpire := l.svcCtx.Config.Jwt.AccessExpire
	accessToken, err := l.getJwtToken(l.svcCtx.Config.Jwt.AccessSecret, now, accessExpire, in.UserId)
	if err != nil {
		return nil, errors.Wrapf(ErrGenerateTokenError, "getJwtToken err: userId = %d, err = %v", in.UserId, err)
	}

	// 将 Token 存入 redis
	userTokenKey := fmt.Sprintf(globalkey.CacheUserTokenKey, in.UserId)
	err = l.svcCtx.RedisClient.Setex(userTokenKey, accessToken, int(accessExpire))
	if err != nil {
		return nil, errors.Wrapf(ErrGenerateTokenError, "Setex err: userId = %d, err = %v", in.UserId, err)
	}
	return &pb.GenerateTokenResp{
		AccessToken:  accessToken,
		AccessExpire: accessExpire,
		RefreshAfter: now + accessExpire/2,
	}, nil
}

func (l *GenerateTokenLogic) getJwtToken(secretKey string, iat, seconds, userId int64) (string, error) {
	return jwt.GetJwtToken(secretKey, iat, seconds, userId)
}
