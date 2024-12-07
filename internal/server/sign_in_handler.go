package server

import (
	"context"
	"fmt"
	"net/http"

	"github.com/opentracing/opentracing-go"
	"github.com/patyukin/mbs-pkg/pkg/errs"
	authpb "github.com/patyukin/mbs-pkg/pkg/proto/auth_v1"
	"github.com/patyukin/mbs-pkg/pkg/proto/error_v1"
)

func (s *Server) SignIn(ctx context.Context, in *authpb.SignInRequest) (*authpb.SignInResponse, error) {
	spanContext := opentracing.SpanFromContext(ctx).Context()
	if spanContext == nil {
		msg := "no span found in context"
		return &authpb.SignInResponse{
			Error: &error_v1.ErrorResponse{
				Code:        http.StatusInternalServerError,
				Message:     msg,
				Description: msg,
			},
		}, nil
	}

	response, err := s.uc.SignInV1UseCase(ctx, in)
	if err != nil {
		return &authpb.SignInResponse{Error: errs.ToErrorResponse(fmt.Errorf("failed s.uc.SignIn: %w", err))}, nil
	}

	return response, nil
}
