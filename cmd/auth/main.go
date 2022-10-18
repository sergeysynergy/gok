package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/caarlos0/env/v6"

	"github.com/sergeysynergy/gok/internal/auth"
	"github.com/sergeysynergy/gok/pkg/utils"
	"github.com/sergeysynergy/gok/tool/conf/service"
	"github.com/sergeysynergy/gok/tool/zapLogger"
)

var (
	buildVersion string
	buildDate    string
)

func main() {
	var err error

	cfg := service.New(
		service.WithDebug(true),
		service.WithDSN("user=gok password=Passw0rd33 host=localhost port=45432 dbname=auth"),
		service.WithAuthAddr(":7000"),
	)

	flag.BoolVar(&cfg.Debug, "debug", cfg.Debug, "run service in debug mode")
	flag.StringVar(&cfg.AuthAddr, "auth", cfg.AuthAddr, "address to start Auth service")
	flag.StringVar(&cfg.DSN, "dsn", cfg.DSN, "Postgres DSN")
	flag.StringVar(&cfg.TrustedSubnet, "ts", cfg.TrustedSubnet, "CIDR string")
	flag.Parse()

	// Rewriting config with environment variables: highest priority.
	err = env.Parse(cfg)
	if err != nil {
		log.Fatalln(err)
	}

	lg := zapLogger.NewServerLogger(cfg.Debug)

	// Warnings for using DEBUG mode.
	if cfg.Debug {
		lg.Warn("ATTENTION: service is running in debug mode")
		lg.Debug(fmt.Sprintf("%#v", cfg))
	}

	// Print build variables that has been set on linking stage, for example:
	// go run -ldflags "-X main.Version=v1.0.1" main.go
	lg.Info(fmt.Sprintf("GoK service version: %s", utils.CheckNA(buildVersion)))
	lg.Info(fmt.Sprintf("GoK build date: %s", utils.CheckNA(buildDate)))

	server := auth.New(cfg, lg)
	server.Run()
}
