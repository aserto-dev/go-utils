package grpcclient

type Config struct {
	Address          string            `json:"address"`
	CACertPath       string            `json:"ca_cert_path"`
	ClientCertPath   string            `json:"client_cert_path"`
	ClientKeyPath    string            `json:"client_key_path"`
	APIKey           string            `json:"api_key"`
	Insecure         bool              `json:"insecure"`
	TimeoutInSeconds int               `json:"timeout_in_seconds"`
	Token            string            `json:"token"`
	Headers          map[string]string `json:"headers"`
}
