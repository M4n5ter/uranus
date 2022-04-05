// Code generated by goctl. DO NOT EDIT.
package types

type FlightOrder struct {
	Sn               string `json:"sn"`               // 订单号
	UserId           int64  `json:"userId"`           // 用户id
	TicketId         int64  `json:"ticketId"`         // 票id
	DepartPosition   string `json:"departPosition"`   // 出发地点
	DepartTime       string `json:"departTime"`       // 出发时间
	ArrivePosisition string `json:"arrivePosisition"` // 降落地点
	ArriveTime       string `json:"arriveTime"`       // 降落地点
	TicketPrice      int64  `json:"ticketPrice"`      // 票价(分)
	Discount         int64  `json:"discount"`         // 折扣(-n%)
	TradeState       int64  `json:"tradeState"`       // 交易状态（-1: 已取消 0:待支付 1:未使用 2:已使用  3:已退款 4:已过期）
	TradeCode        string `json:"tradeCode"`
	OrderTotalPrice  int64  `json:"orderTotalPrice"` // 订单总价(分)(ticketPrice-ticketPrice*discount)
	CreateTime       string `json:"createTime"`      // 订单创建时间
}

type UserFlightOrderListReq struct {
	LastId      int64 `json:"lastId"`      // 最大的 id ，小于等于 0 表示不做最大 id 限制
	PageSize    int64 `json:"pageSize"`    // 查询的条数
	TraderState int64 `json:"traderState"` // 交易状态（-1: 已取消 0:待支付 1:未使用 2:已使用  3:已退款 4:已过期）
}

type UserFlightOrderListResp struct {
	FlightOrderList []string `json:"flightOrderList"` // 用户拥有的订单号
}

type UserFlightOrderDetailReq struct {
	Sn string `json:"sn"`
}

type UserFlightOrderDetailResp struct {
	FlightOrderDetail FlightOrder `json:"flightOrderDetail"`
}
