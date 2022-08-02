package debug

type Config struct {
	Enabled         bool   `json:"enabled"`
	ListenAddress   string `json:"listen_address"`
	ShutdownTimeout int    `json:"shutdown_timeout"`
}
