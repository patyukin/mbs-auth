package server

import (
	"context"
	"fmt"
	"net/http"

	authpb "github.com/patyukin/mbs-pkg/pkg/proto/auth_v1"
	"github.com/patyukin/mbs-pkg/pkg/proto/error_v1"
	"github.com/rs/zerolog/log"
)

func (s *Server) RefreshToken(ctx context.Context, in *authpb.RefreshTokenRequest) (*authpb.RefreshTokenResponse, error) {
	response, err := s.uc.RefreshTokenV1UseCase(ctx, in)
	if err != nil {
		log.Error().Msgf("failed uc.RefreshToken: %v", err)
		return &authpb.RefreshTokenResponse{
			Error: &error_v1.ErrorResponse{
				Code:        http.StatusInternalServerError,
				Message:     "failed RefreshToken",
				Description: fmt.Sprintf("failed uc.RefreshToken: %v", err.Error()),
			},
		}, nil
	}

	return response, nil
}
