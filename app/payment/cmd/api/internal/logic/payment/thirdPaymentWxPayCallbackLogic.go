package payment

import (
	"context"
	"github.com/pkg/errors"
	"github.com/wechatpay-apiv3/wechatpay-go/core/auth/verifiers"
	"github.com/wechatpay-apiv3/wechatpay-go/core/downloader"
	"github.com/wechatpay-apiv3/wechatpay-go/core/notify"
	"github.com/wechatpay-apiv3/wechatpay-go/services/payments"
	"net/http"
	"uranus/app/payment/cmd/rpc/payment"
	"uranus/app/payment/model"
	"uranus/common/xerr"

	"uranus/app/payment/cmd/api/internal/svc"
	"uranus/app/payment/cmd/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

const (
	SUCCESS    = "SUCCESS"    //支付成功
	REFUND     = "REFUND"     //转入退款
	NOTPAY     = "NOTPAY"     //未支付
	CLOSED     = "CLOSED"     //已关闭
	REVOKED    = "REVOKED"    //已撤销（付款码支付）
	USERPAYING = "USERPAYING" //用户支付中（付款码支付）
	PAYERROR   = "PAYERROR"   //支付失败(其他原因，如银行返回失败)
)

var ErrWxPayCallbackError = xerr.NewErrMsg("微信支付回调失败")

type ThirdPaymentWxPayCallbackLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewThirdPaymentWxPayCallbackLogic(ctx context.Context, svcCtx *svc.ServiceContext) ThirdPaymentWxPayCallbackLogic {
	return ThirdPaymentWxPayCallbackLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ThirdPaymentWxPayCallbackLogic) ThirdPaymentWxPayCallback(rw http.ResponseWriter, req *http.Request) (resp *types.ThirdPaymentWxPayCallbackResp, err error) {

	// 读取商户私钥
	_, err = svc.NewWxPayClientV3(l.svcCtx.Config)
	if err != nil {
		return nil, err
	}

	// 获取平台证书访问器
	certVisitor := downloader.MgrInstance().GetCertificateVisitor(l.svcCtx.Config.WxPayConf.MchId)
	handler := notify.NewNotifyHandler(l.svcCtx.Config.WxPayConf.APIv3Key, verifiers.NewSHA256WithRSAVerifier(certVisitor))

	// 校验签名，解析数据
	transaction := &payments.Transaction{}
	_, err = handler.ParseNotifyRequest(context.Background(), req, transaction)
	if err != nil {
		return nil, errors.Wrapf(xerr.NewErrMsg("解析数据失败"), "err: %v", err)
	}

	returnCode := SUCCESS

	err = l.validAndUpdateState(transaction)
	if err != nil {
		returnCode = PAYERROR
	}

	return &types.ThirdPaymentWxPayCallbackResp{ReturnCode: returnCode}, err
}

// 校验和更新相关流水信息
func (l *ThirdPaymentWxPayCallbackLogic) validAndUpdateState(notifyTransaction *payments.Transaction) error {
	// 根据订单号获取支付流水详情
	paymentDetailResp, err := l.svcCtx.PaymentClient.GetPaymentBySn(l.ctx, &payment.GetPaymentBySnReq{Sn: *notifyTransaction.OutTradeNo})
	if err != nil {
		return errors.Wrapf(ErrWxPayCallbackError, "获取支付详情失败, err: %v, notifyTransaction: %+v", err, notifyTransaction)
	}

	// 比对支付金额
	notifyPaytotal := *notifyTransaction.Amount.PayerTotal
	if notifyPaytotal != paymentDetailResp.PaymentDetail.PayTotal {
		return errors.Wrapf(ErrWxPayCallbackError, "支付金额不一致, notifyPaytotal: %v, paymentDetailResp.PaymentDetail.PayTotal: %v, notifyTransaction: %+v", notifyPaytotal, paymentDetailResp.PaymentDetail.PayTotal, notifyTransaction)
	}

	// 判断状态
	payStatus := l.getPayStatusByWXPayTradeState(*notifyTransaction.TradeState)
	if payStatus == model.PaymentPayTradeStateSuccess {
		//  支付通知

		// 更新状态
		if _, err = l.svcCtx.PaymentClient.UpdateTradeState(l.ctx, &payment.UpdateTradeStateReq{
			Sn:             *notifyTransaction.OutTradeNo,
			TradeState:     *notifyTransaction.TradeState,
			TransactionId:  *notifyTransaction.TransactionId,
			TradeType:      *notifyTransaction.TradeType,
			TradeStateDesc: *notifyTransaction.TradeStateDesc,
			PayStatus:      payStatus,
		}); err != nil {
			return errors.Wrapf(ErrWxPayCallbackError, "更新流水状态失败, err: %v, notifyTransaction: %+v", err, notifyTransaction)
		}

	} else if payStatus == model.PaymentPayTradeStateRefund {
		// 退款通知 todo 暂时不需要
	}

	return nil
}

func (l *ThirdPaymentWxPayCallbackLogic) getPayStatusByWXPayTradeState(wxPayTradeState string) int64 {

	switch wxPayTradeState {
	case SUCCESS: //支付成功
		return model.PaymentPayTradeStateSuccess
	case USERPAYING: //支付中
		return model.PaymentPayTradeStateWait
	case REFUND: //已退款
		return model.PaymentPayTradeStateRefund
	default:
		return model.PaymentPayTradeStateFAIL
	}
}
