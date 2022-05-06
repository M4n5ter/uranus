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

type UploadAvatarLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUploadAvatarLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UploadAvatarLogic {
	return &UploadAvatarLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UploadAvatarLogic) UploadAvatar(in *pb.UploadAvatarReq) (*pb.UploadAvatarResp, error) {
	user, err := l.svcCtx.UserModel.FindOne(in.UserId)
	if err != nil && err != userCenterModel.ErrNotFound {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DB_ERROR), "DBERR: %+v", err)
	}

	if user == nil {
		return nil, errors.Wrapf(xerr.NewErrMsg("用户不存在"), "用户不存在,userID: %d", in.UserId)
	}

	user.Avatar = in.Avatar
	err = l.svcCtx.UserModel.UpdateWithVersion(nil, user)
	if err != nil {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DB_ERROR), "DBERR: %+v", err)
	}

	return &pb.UploadAvatarResp{}, nil
}
