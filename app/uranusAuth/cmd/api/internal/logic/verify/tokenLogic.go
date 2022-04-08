package verify

import (
	"context"
	"encoding/json"
	"github.com/golang-jwt/jwt/v4"
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/rest/token"
	"net/http"
	"strings"
	"uranus/app/uranusAuth/cmd/rpc/auth"
	"uranus/common/ctxdata"
	"uranus/common/xerr"

	"uranus/app/uranusAuth/cmd/api/internal/svc"
	"uranus/app/uranusAuth/cmd/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

var ValidateTokenError = xerr.NewErrCode(xerr.TOKEN_EXPIRE_ERROR)

type TokenLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewTokenLogic(ctx context.Context, svcCtx *svc.ServiceContext) TokenLogic {
	return TokenLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *TokenLogic) Token(req *types.VerifyTokenReq, r *http.Request) (resp *types.VerifyTokenResp, err error) {
	authorization := r.Header.Get("Authorization")
	realRequestPath := r.Header.Get("X-Original-Uri")

	if strings.Contains(realRequestPath, "?") {
		realRequestPath = strings.Split(realRequestPath, "?")[0]
	}

	var resultUserId int64
	if l.NoAuthUrls(realRequestPath) {
		// 不需要认证的界面
		// 如果在不需要认证的界面也给了 token ，那么就验证一下取出 uid，没给就不验证了
		if len(authorization) > 0 {
			userId, err := l.isPass(r)
			if err != nil {
				logx.WithContext(l.ctx).Errorf("authorization: %s, realRequestPath: %s", authorization, realRequestPath)
				return nil, err
			}
			if userId == 0 {
				return nil, errors.Wrapf(ValidateTokenError, "from isPass() , userId == 0, authorization: %s, realRequestPath: %s", authorization, realRequestPath)
			}
			resultUserId = userId
		}

	} else {
		// 需要认证的界面
		userId, err := l.isPass(r)
		if err != nil {
			logx.WithContext(l.ctx).Errorf("authorization: %s, realRequestPath: %s", authorization, realRequestPath)
			return nil, err
		}
		if userId == 0 {
			return nil, errors.Wrapf(ValidateTokenError, "from isPass() , userId == 0, authorization: %s, realRequestPath: %s", authorization, realRequestPath)
		}
		resultUserId = userId
	}
	return &types.VerifyTokenResp{
		UserId: resultUserId,
		Ok:     true,
	}, err
}

func (l *TokenLogic) NoAuthUrls(path string) bool {
	for _, url := range l.svcCtx.Config.NoAuthUrls {
		if url == path {
			return true
		}
	}
	return false
}

// 对当前url的授权是否能够通过
func (l *TokenLogic) isPass(r *http.Request) (int64, error) {

	parser := token.NewTokenParser()
	tok, err := parser.ParseToken(r, l.svcCtx.Config.Jwt.AccessSecret, "")
	if err != nil {
		return 0, errors.Wrapf(ValidateTokenError, "JwtAuthLogic isPass  ParseToken err : %v", err)
	}

	if tok.Valid {
		claims, ok := tok.Claims.(jwt.MapClaims) // 解析token中对内容
		if ok {
			userId, _ := claims[ctxdata.CtxKeyJwtUserId].(json.Number).Int64() // 获取userId 并且到后端redis校验是否过期
			if userId <= 0 {
				return 0, errors.Wrapf(ValidateTokenError, "JwtAuthLogic.isPass invalid userId  tokRaw:%s , tokValid :%v ,userId:%d ", tok.Raw, tok.Valid, userId)
			}
			resp, err := l.svcCtx.AuthRpcClient.ValidateToken(l.ctx, &auth.ValidateTokenReq{
				UserId: userId,
				Token:  tok.Raw,
			})
			if err != nil || !resp.Ok {
				return 0, errors.Wrapf(ValidateTokenError, "JwtAuthLogic.isPass IdentityRpc . ValidateToken err:%v ,resp:%+v , tokRaw:%s , tokValid : %v,userId:%d ", err, resp, tok.Raw, tok.Valid, userId)
			}
			return userId, nil
		} else {
			return 0, errors.Wrapf(ValidateTokenError, "tok.Claims is invalid ,tok.Claims ：%+v , claims : %+v , ok:%v ", tok.Claims, claims, ok)
		}
	}
	return 0, errors.Wrapf(xerr.NewErrMsg("token 无效"), "无效的token: %+v", tok)
}
