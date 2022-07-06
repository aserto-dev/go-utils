package apikey

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type (
	// Maps api key to account id.
	Config map[string]string
)

type Authenticator struct {
	cfg Config
}

func New(cfg Config) *Authenticator {
	return &Authenticator{
		cfg: cfg,
	}
}

func (a *Authenticator) Authenticate(apiKey, accountIDOverride string) (string, error) {
	accountID, ok := a.cfg[apiKey]
	if !ok {
		return "", status.Error(codes.Unauthenticated, "invalid api key")
	}

	if accountIDOverride != "" {
		accountID = accountIDOverride
	}

	return accountID, nil
}
