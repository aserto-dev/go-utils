package grpcutil

import "github.com/grpc-ecosystem/grpc-gateway/runtime"

// IncomingHeaderMatcher is a matcher that makes it so that HTTP clients do not have to prefix
// the header key with Grpc-Metadata-.
// see https://grpc-ecosystem.github.io/grpc-gateway/docs/mapping/customizing_your_gateway/#mapping-from-http-request-headers-to-grpc-client-metadata
func IncomingHeaderMatcher(key string) (string, bool) {
	switch key {
	case string(HeaderAsertoTenantID):
		return key, true
	case string(HeaderAsertoAccountID):
		return key, true
	case string(HeaderAsertoRequestID):
		return key, true
	case string(HeaderAsertoTenantKey):
		return key, true
	default:
		return runtime.DefaultHeaderMatcher(key)
	}
}
