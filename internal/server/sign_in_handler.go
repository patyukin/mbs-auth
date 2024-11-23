package server

import (
	"context"
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
				Code:        500,
				Message:     msg,
				Description: msg,
			},
		}, nil
	}

	response, err := s.uc.SignIn(ctx, in)
	if err != nil {
		return &authpb.SignInResponse{Error: errs.ToErrorResponse(err)}, nil
	}

	return response, nil
}
