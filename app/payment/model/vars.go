package model

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var (
	// PayModeWechatPay 微信支付
	PayModeWechatPay = "WECHAT_PAY"
	// PayModeAliPay 支付宝
	PayModeAliPay = "ALI_PAY"
	// PayModeWalletBalance 钱包余额
	PayModeWalletBalance = "WALLET_BALANCE"
)

var ErrNotFound = sqlx.ErrNotFound
