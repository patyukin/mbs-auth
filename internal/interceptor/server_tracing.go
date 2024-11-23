package interceptor

import (
	"context"
	"fmt"
	"github.com/patyukin/mbs-pkg/pkg/tracing"

	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/uber/jaeger-client-go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func ServerTracingInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, info.FullMethod)
	defer span.Finish()

	spanContext, ok := span.Context().(jaeger.SpanContext)
	if ok {
		ctx = metadata.NewOutgoingContext(ctx, metadata.Pairs(tracing.TraceID, spanContext.TraceID().String()))

		header := metadata.New(map[string]string{tracing.TraceID: spanContext.TraceID().String()})
		err := grpc.SendHeader(ctx, header)
		if err != nil {
			return nil, fmt.Errorf("send header: %w", err)
		}
	}

	res, err := handler(ctx, req)
	if err != nil {
		ext.Error.Set(span, true)
		span.SetTag("err", err.Error())
	}

	return res, err
}
