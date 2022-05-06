package logic

import (
	"context"
	"github.com/pkg/errors"
	userCenterModel "uranus/app/userCenter/model"
	"uranus/common/xerr"

	"uranus/app/userCenter/cmd/rpc/internal/svc"
	"uranus/app/userCenter/cmd/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateAvatarLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateAvatarLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateAvatarLogic {
	return &UpdateAvatarLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateAvatarLogic) UpdateAvatar(in *pb.UpdateAvatarReq) (*pb.UpdateAvatarResp, error) {
	user, err := l.svcCtx.UserModel.FindOne(in.UserId)
	if err != nil && err != userCenterModel.ErrNotFound {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DB_ERROR), "DBERR: %+v", err)
	}

	if user == nil {
		return nil, errors.Wrapf(xerr.NewErrMsg("用户不存在"), "用户不存在,userID: %d", in.UserId)
	}

	// 删除原头像 todo

	//更新头像
	user.Avatar = in.Avatar
	err = l.svcCtx.UserModel.UpdateWithVersion(nil, user)
	if err != nil {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DB_ERROR), "DBERR: %+v", err)
	}

	return &pb.UpdateAvatarResp{}, nil
}
