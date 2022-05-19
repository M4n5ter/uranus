package flightInquiry

import (
	"net/http"

	"uranus/common/result"

	"github.com/zeromicro/go-zero/rest/httpx"
	"uranus/app/flightInquiry/cmd/api/internal/logic/flightInquiry"
	"uranus/app/flightInquiry/cmd/api/internal/svc"
	"uranus/app/flightInquiry/cmd/api/internal/types"
)

func GetDiscountFlightsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetDiscountFlightsReq
		if err := httpx.Parse(r, &req); err != nil {
			result.ParamErrorResult(r, w, err)
			return
		}

		l := flightInquiry.NewGetDiscountFlightsLogic(r.Context(), svcCtx)
		resp, err := l.GetDiscountFlights(&req)
		result.HttpResult(r, w, resp, err)
	}
}
