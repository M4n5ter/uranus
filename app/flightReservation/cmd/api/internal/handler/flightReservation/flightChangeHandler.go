package flightReservation

import (
	"net/http"

	"uranus/common/result"

	"github.com/zeromicro/go-zero/rest/httpx"
	"uranus/app/flightReservation/cmd/api/internal/logic/flightReservation"
	"uranus/app/flightReservation/cmd/api/internal/svc"
	"uranus/app/flightReservation/cmd/api/internal/types"
)

func FlightChangeHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.FlightChangeReq
		if err := httpx.Parse(r, &req); err != nil {
			result.ParamErrorResult(r, w, err)
			return
		}

		l := flightReservation.NewFlightChangeLogic(r.Context(), svcCtx)
		resp, err := l.FlightChange(&req)
		result.HttpResult(r, w, resp, err)
	}
}
