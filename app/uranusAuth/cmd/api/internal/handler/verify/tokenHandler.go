package verify

import (
	"fmt"
	"github.com/pkg/errors"
	"net/http"
	"uranus/common/xerr"

	"uranus/app/uranusAuth/cmd/api/internal/logic/verify"
	"uranus/app/uranusAuth/cmd/api/internal/svc"
	"uranus/app/uranusAuth/cmd/api/internal/types"
	"uranus/common/result"

	"github.com/zeromicro/go-zero/rest/httpx"
)

var ErrTokenExpireError = xerr.NewErrCode(xerr.TOKEN_EXPIRE_ERROR)

func TokenHandler(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.VerifyTokenReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}

		l := verify.NewTokenLogic(r.Context(), ctx)
		resp, err := l.Token(&req, r)
		if err == nil && (resp == nil || !resp.Ok) {
			err = errors.Wrapf(ErrTokenExpireError, "jwtAuthHandler JWT Auth no err , userId is zero , req:%+v,resp:%+v", req, resp)
		}

		XUser := "0"
		if resp != nil {
			XUser = fmt.Sprintf("%d", resp.UserId)
		}
		w.Header().Set("x-user", XUser)
		result.HttpResult(r, w, resp, err)
	}
}
