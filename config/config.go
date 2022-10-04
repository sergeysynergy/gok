package config

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
	Addr          string        `env:"ADDRESS" json:"address"`
	StoreFile     string        `env:"STORE_FILE" json:"store_file"`
	Restore       bool          `env:"RESTORE" json:"restore"`
	StoreInterval time.Duration `env:"STORE_INTERVAL"`
	DatabaseDSN   string        `env:"DATABASE_DSN" json:"database_dsn"`
	CryptoKey     string        `env:"CRYPTO_KEY" json:"crypto_key"`
	Key           string        `env:"KEY"`
	ConfigFile    string
}

func NewServerConf() *ServerConf {
	defaultCfg := &ServerConf{
		Addr:          "127.0.0.1:8080",
		StoreFile:     "/tmp/devops-metrics-pgsql.json",
		Restore:       true,
		StoreInterval: 300 * time.Second,
	}

	return defaultCfg
}
