package interceptor

import (
	"context"
	"github.com/rs/zerolog/log"
	"time"

	"google.golang.org/grpc"
)

func LogInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	now := time.Now()

	res, err := handler(ctx, req)
	if err != nil {
		log.Error().Msgf("error: %v, method: %s, req: %v", err, info.FullMethod, req)
	}

	log.Info().Msgf("request, method: %s, req: %v, res: %v, duration: %v", info.FullMethod, req, res, time.Since(now))

	return res, err
}
