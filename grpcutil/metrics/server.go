package metrics

import (
	"net/http"

	"github.com/aserto-dev/go-utils/logger"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog"
	"go.opencensus.io/zpages"
	"google.golang.org/grpc"
)

type Server struct {
	http *http.Server
	cfg  *Config
}

// newMetricsServer creates an http.Server that serves diagnostic metrics.
func NewServer(cfg *Config, log *zerolog.Logger) *Server {
	debugMux := http.NewServeMux()
	zpages.Handle(debugMux, "/debug")

	mux := http.NewServeMux()
	mux.Handle("/debug/", debugMux)
	mux.Handle("/metrics", promhttp.Handler())

	newLogger := log.With().Str("source", "metrics").Logger()
	httpServer := &http.Server{
		ErrorLog: logger.NewSTDLogger(&newLogger),
		Addr:     cfg.ListenAddress,
		Handler:  mux,
	}

	return &Server{http: httpServer, cfg: cfg}
}

func (s *Server) HTTP() *http.Server {
	return s.http
}

// RegisterPrometheusIfEnabled registers prometheus metrics if they are enabled.
func RegisterPrometheusIfEnabled(cfg *Config, srv *grpc.Server) {
	if cfg.GRPC.Counters {
		grpc_prometheus.Register(srv)
	}

	if cfg.GRPC.Durations {
		grpc_prometheus.EnableHandlingTimeHistogram()
	}
}
