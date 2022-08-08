package metrics

type GRPCConfig struct {
	Counters  bool `json:"counters"`
	Durations bool `json:"durations"`
	Gateway   bool `json:"gateway"`
}

func (p GRPCConfig) AllDisabled() bool {
	return !(p.Counters || p.Durations || p.Gateway)
}

// Config defined configuration format for diagnostics and performance metrics.
type Config struct {
	ListenAddress string `json:"listen_address"`

	ZPages bool       `json:"zpages"`
	GRPC   GRPCConfig `json:"grpc"`
	HTTP   bool       `json:"http"`
	DB     bool       `json:"db"`
}

// AllDisabled returns true if all metrics options are disabled. Otherwise, false.
func (m Config) AllDisabled() bool {
	return m.GRPC.AllDisabled() && !m.ZPages && !m.HTTP
}
