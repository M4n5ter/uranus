package payment

import (
	"net/http"

	"uranus/app/payment/cmd/api/internal/logic/payment"
	"uranus/app/payment/cmd/api/internal/svc"
	"uranus/app/payment/cmd/api/internal/types"
	"uranus/common/result"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func LocalPayHandler(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.LocalPaymentReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}

		l := payment.NewLocalPayLogic(r.Context(), ctx)
		resp, err := l.LocalPay(&req)
		result.HttpResult(r, w, resp, err)
	}
}
