package flightOrder

import (
	"net/http"

	"uranus/common/result"

	"github.com/zeromicro/go-zero/rest/httpx"
	"uranus/app/order/cmd/api/internal/logic/flightOrder"
	"uranus/app/order/cmd/api/internal/svc"
	"uranus/app/order/cmd/api/internal/types"
)

func UserFlightOrderDetailHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.UserFlightOrderDetailReq
		if err := httpx.Parse(r, &req); err != nil {
			result.ParamErrorResult(r, w, err)
			return
		}

		l := flightOrder.NewUserFlightOrderDetailLogic(r.Context(), svcCtx)
		resp, err := l.UserFlightOrderDetail(&req)
		result.HttpResult(r, w, resp, err)
	}
}
