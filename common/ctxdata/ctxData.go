package ctxdata

import (
	"context"
	"encoding/json"
)

//从ctx获取uid
var CtxKeyJwtUserId = "jwtUserId"

func GetUidFromCtx(ctx context.Context) int64 {
	uid, _ := ctx.Value(CtxKeyJwtUserId).(json.Number).Int64()
	return uid
}
