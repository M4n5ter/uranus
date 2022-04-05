package user

import (
	"context"
	"github.com/zeromicro/go-zero/core/logx"
	"uranus/app/userCenter/cmd/api/internal/svc"
	"uranus/app/userCenter/cmd/api/internal/types"

	_ "github.com/go-sql-driver/mysql"
)

//var ErrCasbinErr = xerr.NewErrMsg("Casbin 错误")

//var ErrNotFoundErr = xerr.NewErrMsg("获取用户信息失败")

type SearchUserInfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSearchUserInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) SearchUserInfoLogic {
	return SearchUserInfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SearchUserInfoLogic) SearchUserInfo(req *types.SearchUserInfoReq) (resp *types.SearchUserInfoResp, err error) {
	//e := l.svcCtx.CasbinCachedEnforcer
	//err = e.LoadPolicy()
	//if err != nil {
	//	return nil, errors.Wrapf(ErrCasbinErr, "load policy err: %v", err)
	//}
	//// 当前用户 id
	//currentUserId := strconv.Itoa(int(ctxdata.GetUidFromCtx(l.ctx)))
	//// 检查当前用户是否具有查看其它用户 Mobile 的权限
	//if ok, err := e.Enforce(currentUserId, "dom_root", "allUserInfo", "read"); err != nil {
	//	return nil, errors.Wrapf(ErrCasbinErr, "Err in: e.Enforce(%s, \"dom_root\", \"allUserInfo\", \"read\"), err: %v", currentUserId, err)
	//} else {
	//	// 有查看其它用户 Mobile 的权限，或者查看的人是自己
	//	if ok || ctxdata.GetUidFromCtx(l.ctx) == req.UserId {
	//		user, err := l.svcCtx.UsercenterRpcClient.GetUserInfo(l.ctx, &usercenter.GetUserInfoReq{Id: req.UserId})
	//		if err != nil {
	//			return nil, errors.Wrapf(ErrNotFoundErr, "usercenterRpc.GetUserInfo err: %v, userId: %d", err, req.UserId)
	//		}
	//		_ = copier.Copy(resp.UserInfo, user.User)
	//		return resp, nil
	//	}
	//}
	//// 普通用户不能查看他人 Mobile
	//user, err := l.svcCtx.UsercenterRpcClient.GetUserInfo(l.ctx, &usercenter.GetUserInfoReq{Id: req.UserId})
	//if err != nil {
	//	return nil, errors.Wrapf(ErrNotFoundErr, "usercenterRpc.GetUserInfo err: %v, userId: %d", err, req.UserId)
	//}
	//user.User.Mobile = "***********"
	//_ = copier.Copy(resp.UserInfo, user.User)
	return
}
