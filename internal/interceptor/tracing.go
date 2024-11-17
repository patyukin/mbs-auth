package interceptor

import (
	"context"
	"github.com/rs/zerolog/log"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func TraceIDInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		traceIDStrs := md.Get("x-trace-id")
		if len(traceIDStrs) > 0 {
			traceIDStr := traceIDStrs[0]

			traceID, err := trace.TraceIDFromHex(traceIDStr)
			if err == nil {
				spanContext := trace.NewSpanContext(trace.SpanContextConfig{
					TraceID:    traceID,
					SpanID:     trace.SpanID{}, // Пустой SpanID для корневого span
					TraceFlags: trace.FlagsSampled,
					Remote:     true,
				})
				ctx = trace.ContextWithSpanContext(ctx, spanContext)
			} else {
				log.Info().Msgf("Invalid x-trace-id format: %v", err)
			}
		}
	}

	return handler(ctx, req)
}
