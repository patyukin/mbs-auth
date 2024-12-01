package server

import (
	"context"
	"fmt"

	"github.com/patyukin/mbs-pkg/pkg/errs"
	authpb "github.com/patyukin/mbs-pkg/pkg/proto/auth_v1"
)

func (s *Server) AuthorizeUser(ctx context.Context, in *authpb.AuthorizeUserRequest) (*authpb.AuthorizeUserResponse, error) {
	response, err := s.uc.AuthorizeUserV1UseCase(ctx, in)
	if err != nil || response.GetError() != nil {
		return &authpb.AuthorizeUserResponse{Error: errs.ToErrorResponse(fmt.Errorf("failed s.uc.Authorize: %w", err))}, nil
	}

	return response, nil
}
