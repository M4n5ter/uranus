package uniqueid

import (
	"fmt"
	"time"
	"uranus/common/tool"
)

// SnPrefix 生成sn单号
type SnPrefix string

const (
	SN_PREFIX_FLIGHT_ORDER  SnPrefix = "FLO" //机票订单前缀 order/flight_order
	SN_PREFIX_THIRD_PAYMENT SnPrefix = "PMT" //第三方支付流水记录前缀 payment/payment
)

// GenSn 生成单号
func GenSn(snPrefix SnPrefix) string {
	return fmt.Sprintf("%s%s%s", snPrefix, time.Now().Format("20060102150405"), tool.Krand(8, tool.KC_RAND_KIND_NUM))
}
