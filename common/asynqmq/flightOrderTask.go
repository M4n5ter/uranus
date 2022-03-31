package asynqmq

import (
	"encoding/json"
	"uranus/common/xerr"

	"github.com/hibiken/asynq"
	"github.com/pkg/errors"
)

const (
	TypeFlightOrderCloseDelivery = "flight:order:close"
)

// FlightOrderCloseTaskPayload 延迟关闭航班订单task
type FlightOrderCloseTaskPayload struct {
	Sn string
}

func NewFlightOrderCloseTask(sn string) (*asynq.Task, error) {
	payload, err := json.Marshal(FlightOrderCloseTaskPayload{Sn: sn})
	if err != nil {
		return nil, errors.Wrapf(xerr.NewErrMsg("创建延迟关闭机票订单task到asynq失败"), "【addAsynqTaskMarshaError】err : %v , sn : %s", err, sn)
	}
	return asynq.NewTask(TypeFlightOrderCloseDelivery, payload), nil
}
