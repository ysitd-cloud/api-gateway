package grpc

import (
	"context"

	api "code.ysitd.cloud/api/account"
)

func (s *Server) ValidateUserPassword(ctx context.Context, in *api.ValidateUserRequest) (reply *api.ValidateUserReply, err error) {
	reply, err = s.AccountBackend.ValidateUserPassword(ctx, in)
	if err != nil {
		s.Logger.Error(err)
	}
	return
}

func (s *Server) GetUserInfo(ctx context.Context, in *api.GetUserInfoRequest) (reply *api.GetUserInfoReply, err error) {
	reply, err = s.AccountBackend.GetUserInfo(ctx, in)
	if err != nil {
		s.Logger.Error(err)
	}
	return
}

func (s *Server) GetTokenInfo(ctx context.Context, in *api.GetTokenInfoRequest) (reply *api.GetTokenInfoReply, err error) {
	reply, err = s.AccountBackend.GetTokenInfo(ctx, in)
	if err != nil {
		s.Logger.Error(err)
	}
	return
}
