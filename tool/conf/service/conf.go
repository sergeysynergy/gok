// Package service contains GoK service-specific structure.
package service

type Conf struct {
	// Debug defines is service running in debug mode.
	Debug bool `env:"DEBUG"`
	// Addr contains name:port pair value to run service on.
	Addr string `env:"ADDR"`
	// DSN is credential for database connection.
	DSN           string `env:"DATABASE_DSN"`
	// TrustedSubnet is used for CIDR checks in cross-service communications.
	TrustedSubnet string
}

type option func(conf *Conf)

func New(addr string, opts ...option) *Conf {
	const (
		defaultDebug = false
		defaultTrustedSubnet = "127.0.0.1/24"
	)
	cfg := &Conf{
		Debug: defaultDebug,
		Addr:  addr,
		TrustedSubnet: defaultTrustedSubnet,
	}
	for _, opt := range opts {
		opt(cfg)
	}
	return cfg
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
