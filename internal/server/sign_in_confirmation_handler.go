package server

import (
	"context"
	"fmt"
	authpb "github.com/patyukin/mbs-pkg/pkg/proto/auth_v1"
)

func (s *Server) SignInConfirmation(ctx context.Context, in *authpb.SignInConfirmationRequest) (*authpb.SignInConfirmationResponse, error) {
	response, err := s.uc.SignInConfirmationV1UseCase(ctx, in)
	if err != nil {
		return nil, fmt.Errorf("failed s.uc.SignInVerify: %w", err)
	}

	return response, nil
}
