package user

import (
	"net/http"

	"uranus/app/userCenter/cmd/api/internal/logic/user"
	"uranus/app/userCenter/cmd/api/internal/svc"
	"uranus/app/userCenter/cmd/api/internal/types"
	"uranus/common/result"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func SearchUserInfoHandler(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.SearchUserInfoReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}

		l := user.NewSearchUserInfoLogic(r.Context(), ctx)
		resp, err := l.SearchUserInfo(&req)
		result.HttpResult(r, w, resp, err)
	}
}
