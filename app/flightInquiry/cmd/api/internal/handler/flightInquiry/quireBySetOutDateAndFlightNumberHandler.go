package flightInquiry

import (
	"net/http"

	"uranus/app/flightInquiry/cmd/api/internal/logic/flightInquiry"
	"uranus/app/flightInquiry/cmd/api/internal/svc"
	"uranus/app/flightInquiry/cmd/api/internal/types"
	"uranus/common/result"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func QuireBySetOutDateAndFlightNumberHandler(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.QuireBySetOutDateAndFlightNumberReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}

		l := flightInquiry.NewQuireBySetOutDateAndFlightNumberLogic(r.Context(), ctx)
		resp, err := l.QuireBySetOutDateAndFlightNumber(&req)
		result.HttpResult(r, w, resp, err)
	}
}
