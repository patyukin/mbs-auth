package server

import (
	"context"
	"fmt"
	authpb "github.com/patyukin/mbs-pkg/pkg/proto/auth_v1"
	"github.com/patyukin/mbs-pkg/pkg/proto/error_v1"
)

func (s *Server) SignUp(ctx context.Context, in *authpb.SignUpRequest) (*authpb.SignUpResponse, error) {
	response, err := s.uc.SignUpV1UseCase(ctx, in)
	if err != nil {
		return &authpb.SignUpResponse{
			Error: &error_v1.ErrorResponse{
				Code:        500,
				Message:     "Internal Server Error",
				Description: fmt.Sprintf("failed to s.uc.SignUpV1UseCase: %v", err),
			},
		}, nil
	}

	return response, nil
}
