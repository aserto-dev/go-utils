package auth0

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"time"

	"github.com/aserto-dev/go-utils/cerr"
	"github.com/lestrrat-go/jwx/jwk"
	"github.com/lestrrat-go/jwx/jwt"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// AccountIDResolver finds the AccountID associated with the specified identity.
type AccountIDResolver func(ctx context.Context, identity string) (accountID string, err error)

type Authenticator struct {
	cfg      *Config
	logger   *zerolog.Logger
	resolver AccountIDResolver
}

func New(cfg *Config, logger *zerolog.Logger, resolver AccountIDResolver) *Authenticator {
	return &Authenticator{
		cfg:      cfg,
		logger:   logger,
		resolver: resolver,
	}
}

func (a *Authenticator) Authenticate(ctx context.Context, token string) (string, error) {
	parsedToken, err := a.parseToken(ctx, token)
	if err != nil {
		return "", status.Errorf(codes.Unauthenticated, "invalid authentication token: %v", err)
	}

	audienceOK := false
	for _, audience := range parsedToken.Audience() {
		if audience == a.cfg.Audience {
			audienceOK = true
			break
		}
	}
	if !audienceOK {
		return "", status.Error(codes.Unauthenticated, "invalid authentication token; missing audience")
	}

	accountID, err := a.resolver(ctx, parsedToken.Subject())
	if err != nil {
		return "", cerr.ErrAuthenticationFailed.Err(errors.Wrap(err, "unable to resolve identity"))
	}

	return accountID, nil
}

func (a *Authenticator) parseToken(ctx context.Context, token string) (jwt.Token, error) {
	jwtTemp, err := jwt.ParseString(token, jwt.WithValidate(false))
	if err != nil {
		a.logger.Err(err).Msg("jwt parsing failed")
		return nil, err
	}

	options := []jwt.ParseOption{
		jwt.WithValidate(true),
		jwt.WithAcceptableSkew(time.Duration(a.cfg.AcceptableTimeSkewSeconds) * time.Second),
	}

	jwksURL, err := a.jwksURL(ctx, "https://"+a.cfg.Domain)
	if err != nil {
		a.logger.Debug().Str("issuer", jwtTemp.Issuer()).Msg("token didn't have a JWKS endpoint we could use for verification")
	} else {
		jwkSet, err := jwk.Fetch(ctx, jwksURL.String())
		if err != nil {
			return nil, errors.Wrap(err, "failed to fetch JWK set for validation")
		}

		options = append(options, jwt.WithKeySet(jwkSet))
	}

	jwtToken, err := jwt.ParseString(
		token,
		options...,
	)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse jwt")
	}

	return jwtToken, nil
}

func (a *Authenticator) jwksURL(ctx context.Context, baseURL string) (*url.URL, error) {
	const (
		wellknownConfig = `.well-known/openid-configuration`
		wellknownJWKS   = `.well-known/jwks.json`
	)

	u, err := url.Parse(baseURL)
	if err != nil {
		return nil, err
	}

	if u.Scheme == "" {
		return nil, errors.New("no scheme defined for baseURL")
	}

	u.Path = wellknownConfig

	req, err := http.NewRequestWithContext(ctx, "GET", u.String(), nil)
	if err != nil {
		return nil, err
	}

	client := &http.Client{}

	resp, err := client.Do(req)
	if err == nil {
		defer resp.Body.Close()
		var jwksConfig struct {
			URI string `json:"jwks_uri"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&jwksConfig); err == nil {
			if embeddedURI, err := url.Parse(jwksConfig.URI); err == nil {
				return embeddedURI, nil
			}
		}
	}

	u.Path = wellknownJWKS

	return u, nil
}
