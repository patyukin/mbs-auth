package server

import (
	"context"
	"fmt"
	authpb "github.com/patyukin/mbs-pkg/pkg/proto/auth_v1"
)

func (s *Server) SignInVerify(ctx context.Context, in *authpb.SignInVerifyRequest) (*authpb.SignInVerifyResponse, error) {
	response, err := s.uc.SignInVerify(ctx, in)
	if err != nil {
		return nil, fmt.Errorf("failed s.uc.SignInVerify: %w", err)
	}

	return response, nil
}
