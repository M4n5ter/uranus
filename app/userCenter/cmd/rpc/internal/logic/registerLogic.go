package logic

import (
	"context"
	"github.com/jinzhu/copier"
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"uranus/app/uranusAuth/cmd/rpc/auth"
	"uranus/app/userCenter/model"
	"uranus/common/tool"
	"uranus/common/xerr"

	"uranus/app/userCenter/cmd/rpc/internal/svc"
	"uranus/app/userCenter/cmd/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

var ErrRegisterErr = xerr.NewErrMsg("注册失败")
var ErrExistsErr = xerr.NewErrMsg("注册失败，手机号已经被注册")

type RegisterLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterLogic {
	return &RegisterLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *RegisterLogic) Register(in *pb.RegisterReq) (*pb.RegisterResp, error) {
	user, err := l.svcCtx.UserModel.FindOneByMobile(in.Mobile)
	if err != nil && err != model.ErrNotFound {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DB_ERROR), "DB ERR ! ERR: %v", err)
	}
	if user != nil {
		return nil, errors.Wrapf(ErrExistsErr, "mobile exists err, mobile: %s", in.Mobile)
	}

	var userId int64
	if in.AuthType == model.UserAuthTypeSystem {

		// 系统平台
		if err := l.svcCtx.UserModel.Trans(func(session sqlx.Session) error {
			user := new(model.User)
			if len(in.Password) > 0 {
				user.Password = tool.Md5ByString(in.Password)
			} else {
				return errors.Wrapf(ErrRegisterErr, "register err ,can't find password in input")
			}

			if len(in.Nickname) > 0 {
				user.Nickname = in.Nickname
			} else {
				user.Nickname = tool.Krand(tool.KC_RAND_KIND_ALL, 8)
			}

			if in.Sex == model.Female || in.Sex == model.Male {
				user.Sex = in.Sex
			} else {
				user.Sex = model.Unknown
			}

			user.Mobile = in.Mobile
			retUser, err := l.svcCtx.UserModel.Insert(session, user)
			if err != nil {
				return errors.Wrapf(xerr.NewErrCode(xerr.DB_ERROR), "DB ERR when inserting user: %+v , err: %v", user, err)
			}

			lastId, err := retUser.LastInsertId()
			if err != nil {
				return errors.Wrapf(xerr.NewErrCode(xerr.DB_ERROR), "DB ERR when getting the lastInsertId, retUser: %+v, lastId: %d, err: %v", retUser, lastId, err)
			}
			userId = lastId
			userAuth := new(model.UserAuth)
			userAuth.UserId = userId
			userAuth.AuthType = in.AuthType
			userAuth.AuthKey = in.AuthKey
			_, err = l.svcCtx.UserAuthModel.Insert(session, userAuth)
			if err != nil {
				return errors.Wrapf(xerr.NewErrCode(xerr.DB_ERROR), "DB ERR when inserting userAuth: %+v , err: %v", userAuth, err)
			}

			// 为平台用户分配钱包(微信用户不需要)
			wallet := new(model.Wallet)
			wallet.UserId = userId
			wallet.Money = 0
			_, err = l.svcCtx.WalletModel.Insert(session, wallet)
			if err != nil {
				return errors.Wrapf(xerr.NewErrCode(xerr.DB_ERROR), "DB ERR when inserting userAuth: %+v , err: %v", userAuth, err)
			}

			return nil
		}); err != nil {
			return nil, err
		}

	} else if in.AuthType == model.UserAuthTypeSmallWX {

		// 微信小程序
		if err := l.svcCtx.UserModel.Trans(func(session sqlx.Session) error {
			user := new(model.User)
			user.Mobile = in.Mobile
			user.Sex = in.Sex
			user.Nickname = in.Nickname

			retUser, err := l.svcCtx.UserModel.Insert(session, user)
			if err != nil {
				return errors.Wrapf(xerr.NewErrCode(xerr.DB_ERROR), "DB ERR when inserting user: %+v , err: %v", user, err)
			}
			lastId, err := retUser.LastInsertId()
			if err != nil {
				return errors.Wrapf(xerr.NewErrCode(xerr.DB_ERROR), "DB ERR when getting the lastInsertId, retUser: %+v, lastId: %d, err: %v", retUser, lastId, err)
			}
			userId = lastId

			userAuth := new(model.UserAuth)
			userAuth.AuthKey = in.AuthKey
			userAuth.AuthType = in.AuthType
			userAuth.UserId = userId
			_, err = l.svcCtx.UserAuthModel.Insert(session, userAuth)
			if err != nil {
				return errors.Wrapf(xerr.NewErrCode(xerr.DB_ERROR), "DB ERR when inserting userAuth: %+v , err: %v", userAuth, err)
			}

			return nil
		}); err != nil {
			return nil, err
		}
	} else {
		return nil, errors.Wrapf(ErrRegisterErr, "unkown authType: %s, req: %+v", in.AuthType, in)
	}
	authResp, err := l.svcCtx.AuthRpcClient.GenerateToken(l.ctx, &auth.GenerateTokenReq{UserId: userId})
	if err != nil {
		return nil, err
	}
	resp := &pb.RegisterResp{}
	_ = copier.Copy(resp, authResp)
	return resp, nil
}
