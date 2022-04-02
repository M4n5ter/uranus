package payment

import (
	"net/http"

	"uranus/app/payment/cmd/api/internal/logic/payment"
	"uranus/app/payment/cmd/api/internal/svc"
	"uranus/app/payment/cmd/api/internal/types"
	"uranus/common/result"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func ThirdPaymentwxPayHandler(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ThirdPaymentWxPayReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}

		l := payment.NewThirdPaymentwxPayLogic(r.Context(), ctx)
		resp, err := l.ThirdPaymentwxPay(&req)
		result.HttpResult(r, w, resp, err)
	}
}
