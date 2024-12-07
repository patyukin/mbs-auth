package server

import (
	grpcopentracing "github.com/grpc-ecosystem/go-grpc-middleware/tracing/opentracing"
	"github.com/patyukin/mbs-auth/internal/config"
	"google.golang.org/grpc"
)

func NewGRPCServer(_ *config.Config) *grpc.Server {
	return grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			grpcopentracing.UnaryServerInterceptor(),
		),
	)
}
