// Package service contains GoK service-specific structure.
package service

type Conf struct {
	// Debug defines is service running in debug mode.
	Debug bool `env:"DEBUG"`
	// Addr contains name:port pair value to run service on.
	Addr string `env:"ADDR"`
	// DSN is credential for database connection.
	DSN string `env:"DATABASE_DSN"`
}

type option func(conf *Conf)

func New(opts ...option) *Conf {
	const (
		defaultAddr  = ":7000"
		defaultDebug = false
	)
	cfg := &Conf{
		Debug: defaultDebug,
		Addr:  defaultAddr,
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
