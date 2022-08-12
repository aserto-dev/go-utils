package tracing

import (
	"context"
	"time"

	"github.com/aserto-dev/go-utils/grpcutil"
	public_grpcutil "github.com/aserto-dev/go-utils/grpcutil"
	grpcmiddleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"
)

type TracingMiddleware struct {
	logger *zerolog.Logger
}

func NewTracingMiddleware(logger *zerolog.Logger) *TracingMiddleware {
	return &TracingMiddleware{
		logger: logger,
	}
}

var _ public_grpcutil.Middleware = &TracingMiddleware{}

func (m *TracingMiddleware) Unary() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		method, _ := grpc.Method(ctx)

		apiLogger := m.logger.With().
			Str("method", method).
			Fields(grpcutil.KnownContextValueStrings(ctx)).
			Logger()

		apiLogger.Trace().Interface("request", req).Msg("grpc call start")

		newCtx := apiLogger.WithContext(ctx)

		start := time.Now()
		result, err := handler(newCtx, req)
		apiLogger.Trace().Dur("duration-ms", time.Since(start)).Msg("grpc call end")

		return result, err
	}
}

func (m *TracingMiddleware) Stream() grpc.StreamServerInterceptor {
	return func(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		ctx := stream.Context()
		method, _ := grpc.Method(ctx)

		apiLogger := m.logger.With().Str("method", method).Logger()

		apiLogger.Trace().
			Fields(grpcutil.KnownContextValueStrings(ctx)).
			Msg("grpc stream call")

		newCtx := apiLogger.WithContext(ctx)

		wrapped := grpcmiddleware.WrapServerStream(stream)
		wrapped.WrappedContext = newCtx

		return handler(srv, wrapped)
	}
}
