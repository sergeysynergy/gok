// Package service contains GoK service-specific structure.
package service

// Conf is common config structure for both Auth and Storage services:
// Debug defines is service running in debug mode;
// AuthAddr contains name:port pair value to run service on;
// StorageAddr contains name:port pair value to run service on;
// DSN is credential for database connection;
// TrustedSubnet is used for CIDR checks in cross-service communications.
type Conf struct {
	Debug         bool   `env:"DEBUG"`
	AuthAddr      string `env:"AUTH_ADDR"`
	StorageAddr   string `env:"STORAGE_ADDR"`
	DSN           string `env:"DATABASE_DSN"`
	TrustedSubnet string
}

type option func(conf *Conf)

func New(opts ...option) *Conf {
	const (
		defaultDebug         = false
		defaultTrustedSubnet = "127.0.0.1/24"
	)
	cfg := &Conf{
		Debug:         defaultDebug,
		TrustedSubnet: defaultTrustedSubnet,
	}
	for _, opt := range opts {
		opt(cfg)
	}
	return cfg
}

// WithAuthAddr address for Auth gRPC service.
func WithAuthAddr(addr string) option {
	return func(cfg *Conf) {
		cfg.AuthAddr = addr
	}
}

// WithStorageAddr address for Auth gRPC service.
func WithStorageAddr(addr string) option {
	return func(cfg *Conf) {
		cfg.StorageAddr = addr
	}
}

// WithDSN optionally provide to define DSN value.
func WithDSN(dsn string) option {
	return func(cfg *Conf) {
		if dsn != "" {
			cfg.DSN = dsn
		}
	}
}

// WithDebug optionally provide DEBUG value flag.
func WithDebug(debug bool) option {
	return func(cfg *Conf) {
		cfg.Debug = debug
	}
}
