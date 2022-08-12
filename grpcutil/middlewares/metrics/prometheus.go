package metrics

import (
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"google.golang.org/grpc"

	"github.com/aserto-dev/go-utils/grpcutil"
	config "github.com/aserto-dev/go-utils/grpcutil/metrics"
)

func NewMiddlewares(conf config.Config, middlewares ...grpcutil.Middleware) grpcutil.Middlewares {
	if !(conf.GRPC.Counters || conf.GRPC.Durations) {
		// Don't include grpc middleware if counters and durations are disabled.
		return middlewares
	}

	return append(grpcutil.Middlewares{NewPrometheusMiddleware()}, middlewares...)
}

type PrometheusMiddleware struct{}

func NewPrometheusMiddleware() *PrometheusMiddleware {
	return &PrometheusMiddleware{}
}

func (m *PrometheusMiddleware) Unary() grpc.UnaryServerInterceptor {
	return grpc_prometheus.UnaryServerInterceptor
}

func (m *PrometheusMiddleware) Stream() grpc.StreamServerInterceptor {
	return grpc_prometheus.StreamServerInterceptor
}
