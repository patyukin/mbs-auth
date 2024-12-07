package server

import (
	"context"
	"fmt"
	"net/http"

	"github.com/opentracing/opentracing-go"
	authpb "github.com/patyukin/mbs-pkg/pkg/proto/auth_v1"
	"github.com/patyukin/mbs-pkg/pkg/proto/error_v1"
)

func (s *Server) GetBriefUserByID(ctx context.Context, in *authpb.GetBriefUserByIDRequest) (*authpb.GetBriefUserByIDResponse, error) {
	spanContext := opentracing.SpanFromContext(ctx).Context()
	if spanContext == nil {
		msg := "no span found in context"
		return &authpb.GetBriefUserByIDResponse{
			Error: &error_v1.ErrorResponse{
				Code:        http.StatusInternalServerError,
				Message:     msg,
				Description: msg,
			},
		}, nil
	}

	response, err := s.uc.GetBriefUserByID(ctx, in)
	if err != nil {
		return nil, fmt.Errorf("failed uc.GetBriefUserByIDUseCase: %w", err)
	}

	return response, nil
}
