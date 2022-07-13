package authn

import (
	"strings"

	"github.com/aserto-dev/go-utils/cerr"
)

func ParseAuthHeader(val, expectedScheme string) (string, error) {
	splits := strings.SplitN(val, " ", 2)
	if len(splits) < 2 {
		return "", cerr.ErrAuthenticationFailed.Msg("Bad authorization string")
	}
	if !strings.EqualFold(splits[0], expectedScheme) {
		return "", cerr.ErrAuthenticationFailed.Str("expected-scheme", expectedScheme).Msg("Request unauthenticated with expected scheme")
	}
	return splits[1], nil
}
