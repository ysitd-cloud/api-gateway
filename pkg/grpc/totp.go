package grpc

import (
	"context"

	api "code.ysitd.cloud/api/totp"
)

func (s *Server) IssueKey(ctx context.Context, in *api.IssueKeyRequest) (*api.IssueKeyReply, error) {
	return s.TotpBackend.IssueKey(ctx, in)
}

func (s *Server) Validate(ctx context.Context, in *api.ValidateRequest) (*api.ValidateReply, error) {
	return s.TotpBackend.Validate(ctx, in)
}

func (s *Server) RecoverKey(ctx context.Context, in *api.RecoverRequest) (*api.RecoverReply, error) {
	return s.TotpBackend.RecoverKey(ctx, in)
}

func (s *Server) RemoveKey(ctx context.Context, in *api.RemoveKeyRequest) (*api.RemoveKeyReply, error) {
	return s.TotpBackend.RemoveKey(ctx, in)
}
