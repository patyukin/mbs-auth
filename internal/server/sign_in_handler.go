package server

import (
	"context"
	"fmt"
	"github.com/opentracing/opentracing-go"
	"github.com/patyukin/mbs-pkg/pkg/errs"
	authpb "github.com/patyukin/mbs-pkg/pkg/proto/auth_v1"
)

func (s *Server) SignIn(ctx context.Context, in *authpb.SignInRequest) (*authpb.SignInResponse, error) {
	spanContext := opentracing.SpanFromContext(ctx).Context()
	if spanContext == nil {
		return nil, fmt.Errorf("no span found in context")
	}

	response, err := s.uc.SignIn(ctx, in)
	if err != nil {
		return &authpb.SignInResponse{Error: errs.ToErrorResponse(err)}, nil
	}

	return response, nil
}
