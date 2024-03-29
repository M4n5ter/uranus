package logic

import (
	"context"
	"database/sql"
	"github.com/dtm-labs/dtm/client/dtmgrpc"
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"uranus/app/mqueue/cmd/rpc/mqueue"
	"uranus/app/payment/model"
	"uranus/common/xerr"

	"uranus/app/payment/cmd/rpc/internal/svc"
	"uranus/app/payment/cmd/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateTradeStateRollBackLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateTradeStateRollBackLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateTradeStateRollBackLogic {
	return &UpdateTradeStateRollBackLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// UpdateTradeStateRollBack 回滚更新交易状态
func (l *UpdateTradeStateRollBackLogic) UpdateTradeStateRollBack(in *pb.UpdateTradeStateReq) (*pb.UpdateTradeStateResp, error) {

	// 检查输入合法性
	if len(in.Sn) == 0 || in.PayStatus < -1 || in.PayStatus > 3 || !in.PayTime.IsValid() {
		return nil, status.Error(codes.Aborted, errors.Wrapf(xerr.NewErrMsg("非法输入"), "invalid input : %+v", in).Error())
	}

	//1、流水记录确认.
	payment, err := l.svcCtx.PaymentModel.FindOneBySn(l.ctx, in.Sn)
	if err != nil && err != model.ErrNotFound {
		return nil, status.Error(codes.Internal, errors.Wrapf(ERRDBERR, "DBERR: 确认流水记录失败, err: %v, Sn: %s", err, in.Sn).Error())
	}
	if payment == nil {
		return nil, status.Error(codes.Aborted, errors.Wrapf(xerr.NewErrMsg("流水记录不存在"), "Not Found: Sn: %s", in.Sn).Error())
	}

	barrier, err := dtmgrpc.BarrierFromGrpc(l.ctx)
	db, err := sqlx.NewMysql(l.svcCtx.Config.DB.DataSource).RawDB()
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	if err := barrier.CallWithDB(db, func(tx *sql.Tx) error {
		// 废弃交易
		payment.PayStatus = model.CommonPayDiscard
		err := l.svcCtx.PaymentModel.UpdateWithVersion(l.ctx, sqlx.NewSessionFromTx(tx), payment)
		if err != nil {
			return errors.Wrapf(xerr.NewErrMsg("回滚失败"), "err: %v", err)
		}

		return nil
	}); err != nil {
		logx.Errorf("err: %v \n", err)
		// 回滚失败就重试回滚
		return nil, status.Error(codes.Internal, err.Error())
	}

	// 通知其他服务
	_, err = l.svcCtx.MqueueClient.KqPaymentStatusUpdate(l.ctx, &mqueue.KqPaymentStatusUpdateReq{
		OrderSn:   payment.OrderSn,
		PayStatus: payment.PayStatus,
	})
	if err != nil {
		return nil, status.Error(codes.Internal, errors.Wrapf(xerr.NewErrMsg("支付流水状态变更发送到kq失败"), "支付流水状态变更发送到kq失败").Error())
	}

	return &pb.UpdateTradeStateResp{}, nil
}
