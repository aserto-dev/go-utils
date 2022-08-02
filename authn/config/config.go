package config

import (
	"strings"

	"github.com/aserto-dev/go-utils/authn/apikey"
	"github.com/aserto-dev/go-utils/authn/auth0"
	"github.com/aserto-dev/go-utils/authn/sts"
)

type Config struct {
	APIKeys    apikey.Config `json:"api_keys"`
	Auth0      auth0.Config  `json:"auth0"`
	STSService sts.Config    `json:"sts"`
	Options    CallOptions   `json:"options"`
}

type CallOptions struct {
	Default   Options           `json:"default"`
	Overrides []OptionOverrides `json:"overrides"`
}

type Options struct {
	// Require presence of tenant ID
	NeedsTenant bool `json:"needs_tenant"`

	// API Key for machine-to-machine communication, internal to Aserto
	EnableAPIKey bool `json:"enable_api_key"`
	// Tenant-scoped API key, allowing machine-to-machine communication
	EnableMachineKey bool `json:"enable_machine_key"`
	// Auth0 JWT, with an "aserto.com" audience
	EnableAuth0Token bool `json:"enable_auth0_token"`
	// Dex JWT
	EnableDexToken bool `json:"enable_dex_token"`
	// Allows calls without any form of authentication
	EnableAnonymous bool `json:"enable_anonymous"`
	// Allows calls with container identity authentication
	EnableContainerIdentity bool `json:"enable_container_identity"`
}

type OptionOverrides struct {
	// API paths to override
	Paths []string `json:"paths"`
	// Override options
	Override Options `json:"override"`
}

func (co *CallOptions) ForPath(path string) *Options {
	for _, override := range co.Overrides {
		for _, prefix := range override.Paths {
			if strings.HasPrefix(path, prefix) {
				return &override.Override
			}
		}
	}

	return &co.Default
}
