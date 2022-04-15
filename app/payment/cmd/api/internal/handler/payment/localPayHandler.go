package payment

import (
	"net/http"

	"uranus/common/result"

	"github.com/zeromicro/go-zero/rest/httpx"
	"uranus/app/payment/cmd/api/internal/logic/payment"
	"uranus/app/payment/cmd/api/internal/svc"
	"uranus/app/payment/cmd/api/internal/types"
)

func LocalPayHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.LocalPaymentReq
		if err := httpx.Parse(r, &req); err != nil {
			result.ParamErrorResult(r, w, err)
			return
		}

		l := payment.NewLocalPayLogic(r.Context(), svcCtx)
		resp, err := l.LocalPay(&req)
		result.HttpResult(r, w, resp, err)
	}
}
