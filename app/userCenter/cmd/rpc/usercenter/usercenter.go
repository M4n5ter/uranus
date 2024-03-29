// Code generated by goctl. DO NOT EDIT!
// Source: userCenter.proto

package usercenter

import (
	"context"

	"uranus/app/userCenter/cmd/rpc/pb"

	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
)

type (
	AddMoneyReq              = pb.AddMoneyReq
	AddMoneyResp             = pb.AddMoneyResp
	DeductMoneyReq           = pb.DeductMoneyReq
	DeductMoneyResp          = pb.DeductMoneyResp
	GenerateTokenReq         = pb.GenerateTokenReq
	GenerateTokenResp        = pb.GenerateTokenResp
	GetAvatarSrcReq          = pb.GetAvatarSrcReq
	GetAvatarSrcResp         = pb.GetAvatarSrcResp
	GetUserAuthByAuthKeyReq  = pb.GetUserAuthByAuthKeyReq
	GetUserAuthByAuthKeyResp = pb.GetUserAuthByAuthKeyResp
	GetUserAuthByUserIdReq   = pb.GetUserAuthByUserIdReq
	GetUserAuthByUserIdResp  = pb.GetUserAuthByUserIdResp
	GetUserInfoReq           = pb.GetUserInfoReq
	GetUserInfoResp          = pb.GetUserInfoResp
	GetUserMoneyReq          = pb.GetUserMoneyReq
	GetUserMoneyResp         = pb.GetUserMoneyResp
	GetUserWalletReq         = pb.GetUserWalletReq
	GetUserWalletResp        = pb.GetUserWalletResp
	LoginReq                 = pb.LoginReq
	LoginResp                = pb.LoginResp
	RegisterReq              = pb.RegisterReq
	RegisterResp             = pb.RegisterResp
	UpdateUserWalletReq      = pb.UpdateUserWalletReq
	UpdateUserWalletResp     = pb.UpdateUserWalletResp
	UploadAvatarReq          = pb.UploadAvatarReq
	UploadAvatarResp         = pb.UploadAvatarResp
	User                     = pb.User
	UserAuth                 = pb.UserAuth

	Usercenter interface {
		Login(ctx context.Context, in *LoginReq, opts ...grpc.CallOption) (*LoginResp, error)
		Register(ctx context.Context, in *RegisterReq, opts ...grpc.CallOption) (*RegisterResp, error)
		GenerateToken(ctx context.Context, in *GenerateTokenReq, opts ...grpc.CallOption) (*GenerateTokenResp, error)
		GetUserInfo(ctx context.Context, in *GetUserInfoReq, opts ...grpc.CallOption) (*GetUserInfoResp, error)
		GetUserAuthByUserId(ctx context.Context, in *GetUserAuthByUserIdReq, opts ...grpc.CallOption) (*GetUserAuthByUserIdResp, error)
		GetUserAuthByAuthKey(ctx context.Context, in *GetUserAuthByAuthKeyReq, opts ...grpc.CallOption) (*GetUserAuthByAuthKeyResp, error)
		GetUserMoney(ctx context.Context, in *GetUserMoneyReq, opts ...grpc.CallOption) (*GetUserMoneyResp, error)
		UpdateUserWallet(ctx context.Context, in *UpdateUserWalletReq, opts ...grpc.CallOption) (*UpdateUserWalletResp, error)
		AddMoney(ctx context.Context, in *AddMoneyReq, opts ...grpc.CallOption) (*AddMoneyResp, error)
		AddMoneyRollback(ctx context.Context, in *AddMoneyReq, opts ...grpc.CallOption) (*AddMoneyResp, error)
		DeductMoney(ctx context.Context, in *DeductMoneyReq, opts ...grpc.CallOption) (*DeductMoneyResp, error)
		DeductMontyRollBack(ctx context.Context, in *DeductMoneyReq, opts ...grpc.CallOption) (*DeductMoneyResp, error)
		UploadAvatar(ctx context.Context, in *UploadAvatarReq, opts ...grpc.CallOption) (*UploadAvatarResp, error)
		GetAvatarSrc(ctx context.Context, in *GetAvatarSrcReq, opts ...grpc.CallOption) (*GetAvatarSrcResp, error)
		GetUserWallet(ctx context.Context, in *GetUserWalletReq, opts ...grpc.CallOption) (*GetUserWalletResp, error)
	}

	defaultUsercenter struct {
		cli zrpc.Client
	}
)

func NewUsercenter(cli zrpc.Client) Usercenter {
	return &defaultUsercenter{
		cli: cli,
	}
}

func (m *defaultUsercenter) Login(ctx context.Context, in *LoginReq, opts ...grpc.CallOption) (*LoginResp, error) {
	client := pb.NewUsercenterClient(m.cli.Conn())
	return client.Login(ctx, in, opts...)
}

func (m *defaultUsercenter) Register(ctx context.Context, in *RegisterReq, opts ...grpc.CallOption) (*RegisterResp, error) {
	client := pb.NewUsercenterClient(m.cli.Conn())
	return client.Register(ctx, in, opts...)
}

func (m *defaultUsercenter) GenerateToken(ctx context.Context, in *GenerateTokenReq, opts ...grpc.CallOption) (*GenerateTokenResp, error) {
	client := pb.NewUsercenterClient(m.cli.Conn())
	return client.GenerateToken(ctx, in, opts...)
}

func (m *defaultUsercenter) GetUserInfo(ctx context.Context, in *GetUserInfoReq, opts ...grpc.CallOption) (*GetUserInfoResp, error) {
	client := pb.NewUsercenterClient(m.cli.Conn())
	return client.GetUserInfo(ctx, in, opts...)
}

func (m *defaultUsercenter) GetUserAuthByUserId(ctx context.Context, in *GetUserAuthByUserIdReq, opts ...grpc.CallOption) (*GetUserAuthByUserIdResp, error) {
	client := pb.NewUsercenterClient(m.cli.Conn())
	return client.GetUserAuthByUserId(ctx, in, opts...)
}

func (m *defaultUsercenter) GetUserAuthByAuthKey(ctx context.Context, in *GetUserAuthByAuthKeyReq, opts ...grpc.CallOption) (*GetUserAuthByAuthKeyResp, error) {
	client := pb.NewUsercenterClient(m.cli.Conn())
	return client.GetUserAuthByAuthKey(ctx, in, opts...)
}

func (m *defaultUsercenter) GetUserMoney(ctx context.Context, in *GetUserMoneyReq, opts ...grpc.CallOption) (*GetUserMoneyResp, error) {
	client := pb.NewUsercenterClient(m.cli.Conn())
	return client.GetUserMoney(ctx, in, opts...)
}

func (m *defaultUsercenter) UpdateUserWallet(ctx context.Context, in *UpdateUserWalletReq, opts ...grpc.CallOption) (*UpdateUserWalletResp, error) {
	client := pb.NewUsercenterClient(m.cli.Conn())
	return client.UpdateUserWallet(ctx, in, opts...)
}

func (m *defaultUsercenter) AddMoney(ctx context.Context, in *AddMoneyReq, opts ...grpc.CallOption) (*AddMoneyResp, error) {
	client := pb.NewUsercenterClient(m.cli.Conn())
	return client.AddMoney(ctx, in, opts...)
}

func (m *defaultUsercenter) AddMoneyRollback(ctx context.Context, in *AddMoneyReq, opts ...grpc.CallOption) (*AddMoneyResp, error) {
	client := pb.NewUsercenterClient(m.cli.Conn())
	return client.AddMoneyRollback(ctx, in, opts...)
}

func (m *defaultUsercenter) DeductMoney(ctx context.Context, in *DeductMoneyReq, opts ...grpc.CallOption) (*DeductMoneyResp, error) {
	client := pb.NewUsercenterClient(m.cli.Conn())
	return client.DeductMoney(ctx, in, opts...)
}

func (m *defaultUsercenter) DeductMontyRollBack(ctx context.Context, in *DeductMoneyReq, opts ...grpc.CallOption) (*DeductMoneyResp, error) {
	client := pb.NewUsercenterClient(m.cli.Conn())
	return client.DeductMontyRollBack(ctx, in, opts...)
}

func (m *defaultUsercenter) UploadAvatar(ctx context.Context, in *UploadAvatarReq, opts ...grpc.CallOption) (*UploadAvatarResp, error) {
	client := pb.NewUsercenterClient(m.cli.Conn())
	return client.UploadAvatar(ctx, in, opts...)
}

func (m *defaultUsercenter) GetAvatarSrc(ctx context.Context, in *GetAvatarSrcReq, opts ...grpc.CallOption) (*GetAvatarSrcResp, error) {
	client := pb.NewUsercenterClient(m.cli.Conn())
	return client.GetAvatarSrc(ctx, in, opts...)
}

func (m *defaultUsercenter) GetUserWallet(ctx context.Context, in *GetUserWalletReq, opts ...grpc.CallOption) (*GetUserWalletResp, error) {
	client := pb.NewUsercenterClient(m.cli.Conn())
	return client.GetUserWallet(ctx, in, opts...)
}
