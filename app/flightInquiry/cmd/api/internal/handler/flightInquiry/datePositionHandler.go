package flightInquiry

import (
	"net/http"

	"uranus/app/flightInquiry/cmd/api/internal/logic/flightInquiry"
	"uranus/app/flightInquiry/cmd/api/internal/svc"
	"uranus/app/flightInquiry/cmd/api/internal/types"
	"uranus/common/result"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func DatePositionHandler(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.QuireBySetOutDateStartPositionEndPositionReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}

		l := flightInquiry.NewDatePositionLogic(r.Context(), ctx)
		resp, err := l.DatePosition(&req)
		result.HttpResult(r, w, resp, err)
	}
}
