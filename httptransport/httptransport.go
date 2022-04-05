package httptransport

import (
	"crypto/tls"
	"crypto/x509"
	"net/http"
	"os"

	"github.com/friendsofgo/errors"
	"github.com/rs/zerolog"
)

type Config struct {
	CA       []string
	Insecure bool
}

func TransportWithTrustedCAs(log *zerolog.Logger, config *Config) (*http.Transport, error) {
	if config.Insecure {
		return &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}, nil // nolint:gosec // feature used for debugging
	}
	// Get the SystemCertPool, continue with an empty pool on error
	var rootCAs *x509.CertPool

	rootCAs, err := x509.SystemCertPool()
	if err != nil {
		return nil, errors.Wrap(err, "failed to load system cert pool")
	}

	if rootCAs == nil {
		log.Warn().Err(err).Msg("failed to load system ca certs")
		rootCAs = x509.NewCertPool()
	}

	// Read in the cert files
	for _, localCertFile := range config.CA {
		certs, err := os.ReadFile(localCertFile)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to append %q to RootCAs", localCertFile)
		}

		// Append our cert to the system pool
		if ok := rootCAs.AppendCertsFromPEM(certs); !ok {
			log.Warn().Msg("no certs appended, using system certs only")
		}
	}

	// Trust the augmented cert pool in our client
	conf := &tls.Config{
		RootCAs:    rootCAs,
		MinVersion: tls.VersionTLS12,
	}
	return &http.Transport{TLSClientConfig: conf}, nil
}
