package model

import "github.com/zeromicro/go-zero/core/stores/sqlx"

// 支付方式
var (
	// PayModeWechatPay 微信支付
	PayModeWechatPay = "WECHAT_PAY"
	// PayModeAliPay 支付宝
	PayModeAliPay = "ALI_PAY"
	// PayModeWalletBalance 钱包余额
	PayModeWalletBalance = "WALLET_BALANCE"
)

// 支付状态
var (
	PaymentPayTradeStateFAIL    int64 = -1 //支付失败
	PaymentPayTradeStateWait    int64 = 0  //待支付
	PaymentPayTradeStateSuccess int64 = 1  //支付成功
	PaymentPayTradeStateRefund  int64 = 2  //已退款
)

var ErrNotFound = sqlx.ErrNotFound
