package logic

import (
	"context"
	"github.com/jinzhu/copier"
	"github.com/pkg/errors"
	"google.golang.org/protobuf/types/known/timestamppb"
	"uranus/app/payment/model"
	"uranus/common/xerr"

	"uranus/app/payment/cmd/rpc/internal/svc"
	"uranus/app/payment/cmd/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetPaymentBySnLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetPaymentBySnLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetPaymentBySnLogic {
	return &GetPaymentBySnLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// GetPaymentBySn 根据sn查询流水记录
func (l *GetPaymentBySnLogic) GetPaymentBySn(in *pb.GetPaymentBySnReq) (*pb.GetPaymentBySnResp, error) {
	// 检查输入合法性
	if len(in.Sn) == 0 {
		return nil, errors.Wrapf(xerr.NewErrMsg("非法输入"), "invalid input: %s", in.Sn)
	}

	pd, err := l.svcCtx.PaymentModel.FindOneBySn(in.Sn)
	if err != nil && err != model.ErrNotFound {
		return nil, errors.Wrapf(ERRDBERR, "DBERR: 查询流水记录失败 err: %v, Sn: %s", err, in.Sn)
	}

	var resp pb.PaymentDetail
	if pd != nil {
		_ = copier.Copy(&resp, pd)
		resp.CreateTime = timestamppb.New(pd.CreateTime)
		resp.UpdateTime = timestamppb.New(pd.UpdateTime)
		resp.PayTime = timestamppb.New(pd.PayTime)
	}
	return &pb.GetPaymentBySnResp{PaymentDetail: &resp}, nil
}
