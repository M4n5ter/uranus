package ctxdata

import (
	"context"
	"encoding/json"
)

//从ctx获取uid
var CtxKeyJwtUserId = "jwtUserId"

func GetUidFromCtx(ctx context.Context) int64 {
	if v, ok := ctx.Value(CtxKeyJwtUserId).(json.Number); ok {
		uid, _ := v.Int64()
		return uid
	}
	uid, _ := ctx.Value(CtxKeyJwtUserId).(int64)
	return uid
}
