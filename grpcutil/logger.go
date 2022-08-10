package grpcutil

import (
	"context"

	"github.com/rs/zerolog"
)

var knownValueNames = []CtxKey{
	HeaderAsertoTenantID,
	HeaderAsertoAccountID,
	HeaderAsertoRequestID,
	ValueMachineAccountConnectionID,
}

// CompleteLogger returns a logger that contains the
func CompleteLogger(ctx context.Context, log *zerolog.Logger) *zerolog.Logger {
	values := knownContextValueStrings(ctx)
	completeLogger := log.With().Fields(values).Logger()
	return &completeLogger
}

// knownContextValueStrings is the same as KnownContextValues, but uses string keys (useful for logging)
func knownContextValueStrings(ctx context.Context) map[string]interface{} {
	result := map[string]interface{}{}

	for _, k := range knownValueNames {
		v := extract(ctx, k)
		if v != "" {
			result[string(k)] = v
		}
	}

	return result
}

func extract(ctx context.Context, key CtxKey) string {
	id, ok := ctx.Value(key).(string)
	if !ok {
		return ""
	}

	return id
}
