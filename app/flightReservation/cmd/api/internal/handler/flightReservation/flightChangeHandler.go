package flightReservation

import (
	"net/http"

	"uranus/app/flightReservation/cmd/api/internal/logic/flightReservation"
	"uranus/app/flightReservation/cmd/api/internal/svc"
	"uranus/app/flightReservation/cmd/api/internal/types"
	"uranus/common/result"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func FlightChangeHandler(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.FlightChangeReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}

		l := flightReservation.NewFlightChangeLogic(r.Context(), ctx)
		resp, err := l.FlightChange(&req)
		result.HttpResult(r, w, resp, err)
	}
}