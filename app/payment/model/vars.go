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

// 平台内支付状态
var (
	PaymentLocalPayStatusFAIL    int64 = -1 //支付失败
	PaymentLocalPayStatusWait    int64 = 0  //待支付
	PaymentLocalPayStatusSuccess int64 = 1  //支付成功
	PaymentLocalPayStatusRefund  int64 = 2  //已退款
)

// 第三方支付状态 0:未支付 1:支付成功 -1:支付失败 2:已退款
var (
	ThirdPartyPayTradeStateWait    int64 = 0
	ThirdPartyPayTradeStateSuccess int64 = 1
	ThirdPartyPayTradeStateFAIL    int64 = -1
	ThirdPartyPayTradeStateRefund  int64 = 2
)

// 通用支付状态

var (
	CommonPayFAIL    int64 = -1 // 支付失败
	CommonPaySuccess int64 = 1  // 支付成功
	CommonPayWait    int64 = 0  // 待支付
	CommonPayRefund  int64 = 2  // 已退款
	CommonPayDiscard int64 = -2 // 回滚废弃
)

var ErrNotFound = sqlx.ErrNotFound
