package server

import (
	"context"

	"github.com/patyukin/mbs-pkg/pkg/errs"
	authpb "github.com/patyukin/mbs-pkg/pkg/proto/auth_v1"
)

func (s *Server) SignUp(ctx context.Context, in *authpb.SignUpRequest) (*authpb.SignUpResponse, error) {
	response, err := s.uc.SignUpV1UseCase(ctx, in)
	if err != nil {
		return &authpb.SignUpResponse{Error: errs.ToErrorResponse(err)}, nil
	}

	return response, nil
}
