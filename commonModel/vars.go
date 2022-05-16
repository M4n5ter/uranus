package commonModel

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var ErrNotFound = sqlx.ErrNotFound

type Transfer struct {
	// 表示一次中转航班中的各个航班信息
	F []*FlightInfos
}
