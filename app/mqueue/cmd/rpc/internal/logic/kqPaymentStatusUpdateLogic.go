package logic

import (
	"context"
	"encoding/json"
	"github.com/pkg/errors"
	"uranus/common/kqueue"
	"uranus/common/xerr"

	"uranus/app/mqueue/cmd/rpc/internal/svc"
	"uranus/app/mqueue/cmd/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type KqPaymentStatusUpdateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewKqPaymentStatusUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *KqPaymentStatusUpdateLogic {
	return &KqPaymentStatusUpdateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// KqPaymentStatusUpdate 支付流水状态变更发送到kq
func (l *KqPaymentStatusUpdateLogic) KqPaymentStatusUpdate(in *pb.KqPaymentStatusUpdateReq) (*pb.KqPaymentStatusUpdateResp, error) {
	m := kqueue.PaymentUpdatePayStatusNotifyMessage{
		OrderSn:   in.OrderSn,
		PayStatus: in.PayStatus,
	}

	//2、序列化
	body, err := json.Marshal(m)
	if err != nil {
		return nil, errors.Wrapf(xerr.NewErrMsg("kq kqPaymentStatusUpdateLogic  task marshal error "), "kq kqPaymenStatusUpdateLogic  task marshal error , v : %+v", m)
	}

	if err := l.svcCtx.KqueuePaymentUpdatePayStatusClient.Push(string(body)); err != nil {
		return nil, err
	}

	return &pb.KqPaymentStatusUpdateResp{}, nil
}
