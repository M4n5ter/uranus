package logic

import (
	"context"
	"github.com/hibiken/asynq"
	"github.com/pkg/errors"
	"time"
	"uranus/common/asynqmq"
	"uranus/common/xerr"

	"uranus/app/mqueue/cmd/rpc/internal/svc"
	"uranus/app/mqueue/cmd/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type AqDeferFlightOrderCloseLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAqDeferFlightOrderCloseLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AqDeferFlightOrderCloseLogic {
	return &AqDeferFlightOrderCloseLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// AqDeferFlightOrderClose 添加航班订单延迟关闭到 asynq 队列
func (l *AqDeferFlightOrderCloseLogic) AqDeferFlightOrderClose(in *pb.AqDeferFlightOrderCloseReq) (*pb.AqDeferFlightOrderCloseResp, error) {
	task, err := asynqmq.NewFlightOrderCloseTask(in.Sn)
	if err != nil {
		return nil, err
	}

	_, err = l.svcCtx.AsynqClient.Enqueue(task, asynq.ProcessIn(30*time.Minute))
	if err != nil {
		return nil, errors.Wrapf(xerr.NewErrMsg("添加航班订单到延迟队列失败"), "添加航班订单到延迟队列失败 sn:%s ,err:%v", in.Sn, err)
	}

	return &pb.AqDeferFlightOrderCloseResp{}, nil
}
