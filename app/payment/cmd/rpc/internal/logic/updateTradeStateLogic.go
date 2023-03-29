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

type UpdateTradeStateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateTradeStateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateTradeStateLogic {
	return &UpdateTradeStateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// UpdateTradeState 更新交易状态
func (l *UpdateTradeStateLogic) UpdateTradeState(in *pb.UpdateTradeStateReq) (*pb.UpdateTradeStateResp, error) {
	//// 检查输入合法性
	//if len(in.TradeStateDesc) == 0 || len(in.TradeType) == 0 || len(in.Sn) == 0 || len(in.TransactionId) == 0 ||
	//	in.PayStatus < -1 || in.PayStatus > 2 || !in.PayTime.IsValid() {
	//	return nil, errors.Wrapf(xerr.NewErrMsg("非法输入"), "invalid input : %+v", in)
	//}

	//1、流水记录确认.
	payment, err := l.svcCtx.PaymentModel.FindOneBySn(l.ctx, in.Sn)
	if err != nil && err != model.ErrNotFound {
		return nil, status.Error(codes.Internal, errors.Wrapf(ERRDBERR, "DBERR: 确认流水记录失败, err: %v, Sn: %s", err, in.Sn).Error())
	}
	if payment == nil {
		return nil, status.Error(codes.Aborted, errors.Wrapf(xerr.NewErrMsg("流水记录不存在"), "Not Found: Sn: %s", in.Sn).Error())
	}
	//2、判断状态
	if in.PayStatus == model.CommonPayFAIL || in.PayStatus == model.CommonPaySuccess {
		// 想要修改为支付失败或者支付成功的情况

		if payment.PayStatus != model.CommonPayWait {
			return nil, status.Error(codes.Aborted, errors.Wrapf(xerr.NewErrMsg("只有待支付的订单可以修改为支付失败状态或者支付成功状态"), "当前流水状态非待支付状态，不可修改为支付成功或失败, in: %+v", in).Error())
		}

	} else if in.PayStatus == model.CommonPayRefund {
		// 要修改为退款成功的情况
		if payment.PayStatus != model.CommonPaySuccess {
			return nil, status.Error(codes.Aborted, errors.Wrapf(xerr.NewErrMsg("只有付款成功的订单才能退款"), "修改支付流水记录为退款失败，当前支付流水未支付成功无法退款 in : %+v", in).Error())
		}
	} else {
		return nil, status.Error(codes.Aborted, errors.Wrapf(xerr.NewErrMsg("当前不支持该状态"), "不支持的支付流水状态: %+v", in).Error())
	}
	//3、更新.

	// 开一个 barrier
	barrier, err := dtmgrpc.BarrierFromGrpc(l.ctx)
	db, err := sqlx.NewMysql(l.svcCtx.Config.DB.DataSource).RawDB()
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	if err := barrier.CallWithDB(db, func(tx *sql.Tx) error {
		if len(in.TradeState) > 0 {
			// 这种情况下是第三方支付，下方信息都是第三方提供，平台内支付无需提供
			payment.TradeState = in.TradeState
			payment.TransactionId = in.TransactionId
			payment.TradeStateDesc = in.TradeStateDesc
			payment.TradeType = in.TradeType
		}
		payment.PayStatus = in.PayStatus
		payment.PayTime = in.PayTime.AsTime()
		payment.PayMode = in.PayMode
		err = l.svcCtx.PaymentModel.UpdateWithVersion(l.ctx, sqlx.NewSessionFromTx(tx), payment)
		if err != nil {
			return errors.Wrapf(ERRDBERR, "更新支付流水失败: payment: %+v", payment)
		}
		return nil

	}); err != nil {
		return nil, status.Error(codes.Aborted, err.Error())
	}

	//4、通知其他服务
	_, err = l.svcCtx.MqueueClient.KqPaymentStatusUpdate(l.ctx, &mqueue.KqPaymentStatusUpdateReq{
		OrderSn:   payment.OrderSn,
		PayStatus: payment.PayStatus,
	})
	if err != nil {
		return nil, status.Error(codes.Internal, errors.Wrapf(xerr.NewErrMsg("支付流水状态变更发送到kq失败"), "支付流水状态变更发送到kq失败").Error())
	}
	return &pb.UpdateTradeStateResp{}, nil
}
