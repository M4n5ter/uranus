package model

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var ErrNotFound = sqlx.ErrNotFound

// FlightOrderTradeStateCancel 订单取消
var FlightOrderTradeStateCancel int64 = -1

// FlightOrderTradeStateWaitPay 等待支付
var FlightOrderTradeStateWaitPay int64 = 0

// FlightOrderTradeStateWaitUse 等待使用
var FlightOrderTradeStateWaitUse int64 = 1

// FlightOrderTradeStateUsed 已使用
var FlightOrderTradeStateUsed int64 = 2

// FlightOrderTradeStateRefund 退款
var FlightOrderTradeStateRefund int64 = 3

// FlightOrderTradeStateExpire 订单过期
var FlightOrderTradeStateExpire int64 = 4

// FlightOrderTradeStateDiscard 废弃订单
var FlightOrderTradeStateDiscard int64 = -2
