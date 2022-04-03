// Code generated by goctl. DO NOT EDIT.
package handler

import (
	"net/http"

	payment "uranus/app/payment/cmd/api/internal/handler/payment"
	"uranus/app/payment/cmd/api/internal/svc"

	"github.com/zeromicro/go-zero/rest"
)

func RegisterHandlers(server *rest.Server, serverCtx *svc.ServiceContext) {
	server.AddRoutes(
		[]rest.Route{
			{
				Method:  http.MethodPost,
				Path:    "/payment/localpayment",
				Handler: payment.LocalPayHandler(serverCtx),
			},
			{
				Method:  http.MethodPost,
				Path:    "/payment/paymentWxPay",
				Handler: payment.ThirdPaymentwxPayHandler(serverCtx),
			},
			{
				Method:  http.MethodPost,
				Path:    "/payment/paymentWxPayCallback",
				Handler: payment.ThirdPaymentWxPayCallbackHandler(serverCtx),
			},
		},
		rest.WithPrefix("/payment/v1"),
	)
}
