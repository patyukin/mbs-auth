package server

import (
	"context"
	"fmt"
	authpb "github.com/patyukin/mbs-pkg/pkg/proto/auth_v1"
	"github.com/rs/zerolog/log"
)

func (s *Server) RefreshToken(ctx context.Context, in *authpb.RefreshTokenRequest) (*authpb.RefreshTokenResponse, error) {
	response, err := s.uc.RefreshToken(ctx, in)
	if err != nil {
		log.Error().Msgf("failed uc.RefreshToken: %v", err)
		return nil, fmt.Errorf("failed RefreshToken")
	}

	return response, nil
}
