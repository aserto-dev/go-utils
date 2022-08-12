package grpcutil

import (
	"context"

	"github.com/rs/zerolog"
)

// CompleteLogger returns a logger that contains the
func CompleteLogger(ctx context.Context, log *zerolog.Logger) *zerolog.Logger {
	values := KnownContextValueStrings(ctx)
	completeLogger := log.With().Fields(values).Logger()
	return &completeLogger
}
