package verify

import (
	"net/http"

	"uranus/app/uranusAuth/cmd/api/internal/logic/verify"
	"uranus/app/uranusAuth/cmd/api/internal/svc"
	"uranus/app/uranusAuth/cmd/api/internal/types"
	"uranus/common/result"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func TokenHandler(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.VerifyTokenReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}

		l := verify.NewTokenLogic(r.Context(), ctx)
		resp, err := l.Token(&req, r)
		result.HttpResult(r, w, resp, err)
	}
}
