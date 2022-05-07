package logic

import (
	"context"

	"github.com/pkg/errors"
	"github.com/qiniu/go-sdk/v7/storage"

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

	if len(user.Avatar) > 0 {
		// 原本有头像需要删除 oss 中原本的头像

		mac := l.svcCtx.GenQiniuOSSMAC()
		cfg := storage.Config{
			// 是否使用https域名进行资源管理
			UseHTTPS: true,
		}
		// 指定空间所在的区域，如果不指定将自动探测
		// 如果没有特殊需求，默认不需要指定
		//cfg.Zone=&storage.ZoneHuadong
		bucketManager := storage.NewBucketManager(mac, &cfg)

		bucket := l.svcCtx.Config.QiniuOSS.Bucket
		key := user.Avatar
		err := bucketManager.Delete(bucket, key)
		if err != nil {
			logx.Errorf("删除用户原有头像失败,ERR: %+v", err)
		}
	}

	user.Avatar = in.Avatar
	err = l.svcCtx.UserModel.UpdateWithVersion(nil, user)
	if err != nil {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DB_ERROR), "DBERR: %+v", err)
	}

	return &pb.UploadAvatarResp{}, nil
}
