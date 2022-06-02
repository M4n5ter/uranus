package user

import (
	"net/http"

	"uranus/common/result"

	"github.com/zeromicro/go-zero/rest/httpx"
	"uranus/app/userCenter/cmd/api/internal/logic/user"
	"uranus/app/userCenter/cmd/api/internal/svc"
	"uranus/app/userCenter/cmd/api/internal/types"
)

func GetUserWalletHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetUserWalletReq
		if err := httpx.Parse(r, &req); err != nil {
			result.ParamErrorResult(r, w, err)
			return
		}

		l := user.NewGetUserWalletLogic(r.Context(), svcCtx)
		resp, err := l.GetUserWallet(&req)
		result.HttpResult(r, w, resp, err)
	}
}
