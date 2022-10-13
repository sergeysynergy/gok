package serverConf

import (
	"time"
)

type Config struct {
	SyncInterval time.Duration
	//Addr             string        `env:"ADDRESS" json:"address"`
	//MyReportInterval Duration      `json:"report_interval"`
	//MyPollInterval   Duration      `json:"poll_interval"`
	//ReportInterval   time.Duration `env:"REPORT_INTERVAL"`
	//PollInterval     time.Duration `env:"POLL_INTERVAL"`
	//Key              string        `env:"KEY"`
	//CryptoKey        string        `env:"CRYPTO_KEY"`
}

func New() *Config {
	defaultCfg := &Config{}

	return defaultCfg
}

type ServerConf struct {
	Addr string
	DSN  string
	//CryptoKey     string        `env:"CRYPTO_KEY" json:"crypto_key"`
	//Key           string        `env:"KEY"`
}

type ServerOption func(conf *ServerConf)

func NewServerConf(opts ...ServerOption) *ServerConf {
	cfg := &ServerConf{
		Addr: "127.0.0.1:8080",
	}
	for _, opt := range opts {
		opt(cfg)
	}
	return cfg
}

func WithDSN(dsn string) ServerOption {
	return func(cfg *ServerConf) {
		if dsn != "" {
			cfg.DSN = dsn
		}
	}
}
