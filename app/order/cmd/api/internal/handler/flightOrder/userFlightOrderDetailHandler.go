package flightOrder

import (
	"net/http"

	"uranus/app/order/cmd/api/internal/logic/flightOrder"
	"uranus/app/order/cmd/api/internal/svc"
	"uranus/app/order/cmd/api/internal/types"
	"uranus/common/result"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func UserFlightOrderDetailHandler(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.UserFlightOrderDetailReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}

		l := flightOrder.NewUserFlightOrderDetailLogic(r.Context(), ctx)
		resp, err := l.UserFlightOrderDetail(&req)
		result.HttpResult(r, w, resp, err)
	}
}
