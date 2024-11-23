package server

import (
	"context"
	"fmt"
	authpb "github.com/patyukin/mbs-pkg/pkg/proto/auth_v1"
)

func (s *Server) SignUp(ctx context.Context, in *authpb.SignUpRequest) (*authpb.SignUpResponse, error) {
	response, err := s.uc.SignUpV1UseCase(ctx, in)
	if err != nil {
		return nil, fmt.Errorf("failed s.uc.SignIn: %w", err)
	}

	return response, nil
}
