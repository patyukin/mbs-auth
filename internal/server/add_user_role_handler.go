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

func (s *Server) AddUserRole(ctx context.Context, in *authpb.AddUserRoleRequest) (*authpb.AddUserRoleResponse, error) {
	spanContext := opentracing.SpanFromContext(ctx).Context()
	if spanContext == nil {
		return &authpb.AddUserRoleResponse{
			Error: &error_v1.ErrorResponse{
				Code:        http.StatusInternalServerError,
				Message:     "Internal Server Error",
				Description: "no span found in context",
			},
		}, nil
	}

	response, err := s.uc.AddUserRoleV1UseCase(ctx, in)
	if err != nil {
		return &authpb.AddUserRoleResponse{
			Error: errs.ToErrorResponse(fmt.Errorf("failed uc.AddUserRoleV1UseCase: %w", err)),
		}, nil
	}

	return response, nil
}
