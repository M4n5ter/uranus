// Code generated by goctl. DO NOT EDIT.
package handler

import (
	"net/http"

	user "uranus/app/userCenter/cmd/api/internal/handler/user"
	"uranus/app/userCenter/cmd/api/internal/svc"

	"github.com/zeromicro/go-zero/rest"
)

func RegisterHandlers(server *rest.Server, serverCtx *svc.ServiceContext) {
	server.AddRoutes(
		[]rest.Route{
			{
				Method:  http.MethodPost,
				Path:    "/user/register",
				Handler: user.RegisterHandler(serverCtx),
			},
			{
				Method:  http.MethodPost,
				Path:    "/user/login",
				Handler: user.LoginHandler(serverCtx),
			},
		},
		rest.WithPrefix("/usercenter/v1"),
	)

	server.AddRoutes(
		[]rest.Route{
			{
				Method:  http.MethodPost,
				Path:    "/user/detail",
				Handler: user.DetailHandler(serverCtx),
			},
			{
				Method:  http.MethodPost,
				Path:    "/user/searchUserInfo",
				Handler: user.SearchUserInfoHandler(serverCtx),
			},
			{
				Method:  http.MethodPost,
				Path:    "/user/wxMiniAuth",
				Handler: user.WxMiniAuthHandler(serverCtx),
			},
			{
				Method:  http.MethodGet,
				Path:    "/user/avatar/getuptoken",
				Handler: user.GetAvatarUpTokenHandler(serverCtx),
			},
			{
				Method:  http.MethodGet,
				Path:    "/user/avatar/getavatarsrc",
				Handler: user.GetAvatarSrcHandler(serverCtx),
			},
			{
				Method:  http.MethodGet,
				Path:    "/user/avatar/upload",
				Handler: user.UploadAvatarHandler(serverCtx),
			},
		},
		rest.WithJwt(serverCtx.Config.JwtAuth.AccessSecret),
		rest.WithPrefix("/usercenter/v1"),
	)
}
