package payment

import (
	"context"
	"github.com/pkg/errors"
	"github.com/wechatpay-apiv3/wechatpay-go/core"
	"uranus/app/order/cmd/rpc/order"
	"uranus/app/payment/cmd/rpc/payment"
	"uranus/app/payment/model"
	"uranus/app/userCenter/cmd/rpc/usercenter"
	userCenterModel "uranus/app/userCenter/model"
	"uranus/common/ctxdata"
	"uranus/common/xerr"

	"uranus/app/payment/cmd/api/internal/svc"
	"uranus/app/payment/cmd/api/internal/types"

	"github.com/wechatpay-apiv3/wechatpay-go/services/payments/jsapi"
	"github.com/zeromicro/go-zero/core/logx"
)

type ThirdPaymentwxPayLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewThirdPaymentwxPayLogic(ctx context.Context, svcCtx *svc.ServiceContext) ThirdPaymentwxPayLogic {
	return ThirdPaymentwxPayLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ThirdPaymentwxPayLogic) ThirdPaymentwxPay(req *types.ThirdPaymentWxPayReq) (resp *types.ThirdPaymentWxPayResp, err error) {
	var totalPrice int64 // 当前订单支付总金额(分)
	// 获取订单总价
	orderDetail, err := l.svcCtx.OrderClient.FlightOrderDetail(l.ctx, &order.FlightOrderDetailReq{Sn: req.OrderSn})
	if err != nil {
		return nil, err
	}
	if orderDetail == nil {
		return nil, errors.Wrapf(xerr.NewErrMsg("订单不存在"), "orderSn: %s", req.OrderSn)
	}
	totalPrice = orderDetail.FlightOrder.OrderTotalPrice

	// 创建微信预支付订单
	prePayResp, err := l.createWxPrePayOrder(req.OrderSn, totalPrice)
	if err != nil {
		return nil, err
	}
	resp = &types.ThirdPaymentWxPayResp{
		Appid:     l.svcCtx.Config.WxMiniConf.AppId,
		NonceStr:  *prePayResp.NonceStr,
		PaySign:   *prePayResp.PaySign,
		Package:   *prePayResp.Package,
		Timestamp: *prePayResp.TimeStamp,
		SignType:  *prePayResp.SignType,
	}
	return
}

// 获取支付航班当前订单的价格以及信息
func (l *ThirdPaymentwxPayLogic) createWxPrePayOrder(orderSn string, totalPrice int64) (*jsapi.PrepayWithRequestPaymentResponse, error) {

	// 获取用户 openid
	userID := ctxdata.GetUidFromCtx(l.ctx)
	userResp, err := l.svcCtx.UserCenterClient.GetUserAuthByUserId(l.ctx, &usercenter.GetUserAuthByUserIdReq{
		UserId:   userID,
		AuthType: userCenterModel.UserAuthTypeSmallWX,
	})
	if err != nil {
		return nil, errors.Wrapf(xerr.NewErrMsg("获取用户 openid 失败"), "err: %v, userID: %d, orderSn: %s", err, userID, orderSn)
	}
	if userResp.UserAuth == nil {
		return nil, errors.Wrapf(xerr.NewErrMsg("用户 openid 不存在,请先进行微信授权再支付"), "userID: %d, orderSn: %s", userID, orderSn)
	}
	openID := userResp.UserAuth.AuthKey

	// 创建本地支付流水记录
	createPaymentResp, err := l.svcCtx.PaymentClient.CreatePayment(l.ctx, &payment.CreatePaymentReq{
		UserID:   userID,
		PayModel: model.PayModeWechatPay,
		PayTotal: totalPrice,
		OrderSn:  orderSn,
	})
	if err != nil {
		return nil, errors.Wrapf(xerr.NewErrMsg("创建支付流水记录失败"), "err: %v, userID: %d, orderSn: %s, totalPrice: %d", err, userID, orderSn, totalPrice)
	}

	// 创建微信预处理订单
	wxPayClient, err := svc.NewWxPayClientV3(l.svcCtx.Config)
	if err != nil {
		return nil, err
	}

	jsapiSvc := jsapi.JsapiApiService{
		Client: wxPayClient,
	}

	// 得到prepay_id，以及调起支付所需的参数和签名
	resp, _, err := jsapiSvc.PrepayWithRequestPayment(l.ctx, jsapi.PrepayRequest{
		Appid:       core.String(l.svcCtx.Config.WxMiniConf.AppId),
		Mchid:       core.String(l.svcCtx.Config.WxPayConf.MchId),
		Description: core.String("航班支付订单"),
		OutTradeNo:  core.String(createPaymentResp.Sn),
		NotifyUrl:   core.String(l.svcCtx.Config.WxPayConf.NotifyUrl),
		Amount:      &jsapi.Amount{Total: core.Int64(totalPrice)},
		Payer:       &jsapi.Payer{Openid: core.String(openID)},
	})
	if err != nil {
		return nil, errors.Wrapf(xerr.NewErrMsg("创建微信预处理订单失败"), "err: %v, userID: %d, orderSn: %s, totalPrice: %d", err, userID, orderSn, totalPrice)
	}

	return resp, nil
}
