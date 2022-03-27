package flightInquiry

import (
	"net/http"

	"uranus/app/flightInquiry/cmd/api/internal/logic/flightInquiry"
	"uranus/app/flightInquiry/cmd/api/internal/svc"
	"uranus/app/flightInquiry/cmd/api/internal/types"
	"uranus/common/result"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func QuireBySetOutDateStartPositionEndPositionReqHandler(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.QuireBySetOutDateStartPositionEndPositionReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}

		l := flightInquiry.NewQuireBySetOutDateStartPositionEndPositionReqLogic(r.Context(), ctx)
		resp, err := l.QuireBySetOutDateStartPositionEndPositionReq(&req)
		result.HttpResult(r, w, resp, err)
	}
}
