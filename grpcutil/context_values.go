package grpcutil

import (
	"context"

	"github.com/google/uuid"
)

type CtxKey string

var knownValueNames = []CtxKey{
	HeaderAsertoTenantID,
	HeaderAsertoAccountID,
	HeaderAsertoRequestID,
	ValueMachineAccountConnectionID,
}

// If a call is authenticated with a machine account, these hold info about the connection and cert that resolved
// that machine account.
var (
	ValueMachineAccountConnectionID = CtxKey("machine-account-connection-id")
	ValueMachineAccountCertID       = CtxKey("machine-account-cert-id")
)

func IsValidUUID(uid string) bool {
	_, err := uuid.Parse(uid)
	return err == nil
}

// KnownContextValueStrings is the same as KnownContextValues, but uses string keys (useful for logging)
func KnownContextValueStrings(ctx context.Context) map[string]interface{} {
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
