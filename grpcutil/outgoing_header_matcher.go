package grpcutil

import "github.com/grpc-ecosystem/grpc-gateway/runtime"

// OutgoingHeaderMatcher filters headers that are allowed to reach the client
func OutgoingHeaderMatcher(key string) (string, bool) {
	switch key {
	case string(HeaderAsertoRequestIDLowercase):
		return key, true
	default:
		return runtime.DefaultHeaderMatcher(key)
	}
}
