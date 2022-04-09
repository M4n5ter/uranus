package logic

import (
	"context"
	"github.com/pkg/errors"
	"uranus/app/userCenter/model"
	"uranus/common/xerr"

	"uranus/app/userCenter/cmd/rpc/internal/svc"
	"uranus/app/userCenter/cmd/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type AddMoneyLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAddMoneyLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddMoneyLogic {
	return &AddMoneyLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *AddMoneyLogic) AddMoney(in *pb.AddMoneyReq) (*pb.AddMoneyResp, error) {
	// 检查输入
	if in.Money < 0 {
		return nil, errors.Wrapf(xerr.NewErrMsg("增加的钱不能小于 0 ，扣钱请使用隔壁rpc"), "")
	}

	// 检查用户是否存在
	_, err := l.svcCtx.UserModel.FindOne(in.UserId)
	if err != nil {
		if err == model.ErrNotFound {
			return nil, errors.Wrapf(xerr.NewErrMsg("用户不存在"), "Not Found userID: %d", in.UserId)
		}
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DB_ERROR), "DBERR: %v", err)
	}
	// 检查用户是否有钱包
	wallet, err := l.svcCtx.WalletModel.FindOneByUserId(in.UserId)
	if err != nil && err != model.ErrNotFound {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DB_ERROR), "DBERR: %v", err)
	}

	if wallet == nil {
		return nil, errors.Wrapf(xerr.NewErrMsg("该用户没有钱包"), "userID: %d", in.UserId)
	}
	// 加钱
	wallet.Money = wallet.Money + in.Money
	err = l.svcCtx.WalletModel.UpdateWithVersion(nil, wallet)
	if err != nil {
		return nil, errors.Wrapf(xerr.NewErrMsg("加钱失败"), "DBERR: %v", err)
	}

	return &pb.AddMoneyResp{}, nil
}
