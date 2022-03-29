package user

import (
	"net/http"

	"uranus/app/userCenter/cmd/api/internal/logic/user"
	"uranus/app/userCenter/cmd/api/internal/svc"
	"uranus/app/userCenter/cmd/api/internal/types"
	"uranus/common/result"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func WxMiniAuthHandler(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.WXMiniAuthReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}

		l := user.NewWxMiniAuthLogic(r.Context(), ctx)
		resp, err := l.WxMiniAuth(&req)
		result.HttpResult(r, w, resp, err)
	}
}
