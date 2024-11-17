package server

import (
	"context"
	"fmt"
	"github.com/opentracing/opentracing-go"
	authpb "github.com/patyukin/mbs-pkg/pkg/proto/auth_v1"
	"github.com/rs/zerolog/log"
)

func (s *Server) AddUserRole(ctx context.Context, in *authpb.AddUserRoleRequest) (*authpb.AddUserRoleResponse, error) {
	spanContext := opentracing.SpanFromContext(ctx).Context()
	if spanContext == nil {
		log.Fatal().Msgf("no span found in context")
	}

	response, err := s.uc.AddUserRole(ctx, in)
	if err != nil {
		return nil, fmt.Errorf("failed s.uc.AddUserRole: %w", err)
	}

	return response, nil
}
