package flightInquiry

import (
	"net/http"

	"uranus/common/result"

	"github.com/zeromicro/go-zero/rest/httpx"
	"uranus/app/flightInquiry/cmd/api/internal/logic/flightInquiry"
	"uranus/app/flightInquiry/cmd/api/internal/svc"
	"uranus/app/flightInquiry/cmd/api/internal/types"
)

func GetFlightsByPriceRangeHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetFlightsByPriceRangeReq
		if err := httpx.Parse(r, &req); err != nil {
			result.ParamErrorResult(r, w, err)
			return
		}

		l := flightInquiry.NewGetFlightsByPriceRangeLogic(r.Context(), svcCtx)
		resp, err := l.GetFlightsByPriceRange(&req)
		result.HttpResult(r, w, resp, err)
	}
}
