package server

import (
	"context"
	"fmt"

	"github.com/patyukin/mbs-pkg/pkg/errs"
	authpb "github.com/patyukin/mbs-pkg/pkg/proto/auth_v1"
)

func (s *Server) SignInConfirmation(ctx context.Context, in *authpb.SignInConfirmationRequest) (*authpb.SignInConfirmationResponse, error) {
	response, err := s.uc.SignInConfirmationV1UseCase(ctx, in)
	if err != nil {
		return &authpb.SignInConfirmationResponse{
			Error: errs.ToErrorResponse(fmt.Errorf("failed s.uc.SignInConfirmationV1UseCase: %w", err)),
		}, nil
	}

	return response, nil
}
