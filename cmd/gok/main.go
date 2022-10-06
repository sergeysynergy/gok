package main

import (
	"fmt"
	"log"
	"os"

	"github.com/google/uuid"
)

var (
	buildVersion string
	buildDate    string
)

func main() {
	// Выведем номер версии, сборки и комит, если доступны.
	// Для задания переменных рекомендуется использовать опции линковщика, например:
	// go run -ldflags "-X main.buildVersion=v1.0.1" main.go
	//fmt.Printf("Build version: %s\n", utils.CheckNA(buildVersion))
	//fmt.Printf("Build date: %s\n", utils.CheckNA(buildDate))

	// TODO: uncomment
	//if len(os.Args) < 2 {
	//	help()
	//	return
	//}

	// TODO: Add logger init here

	var home string
	for _, v := range os.Args {
		switch v {
		case "--home":
			home = "test"
		}
	}

	home = getHome(home)
	log.Println("[DEBUG] home:", home)

	id := uuid.New()
	fmt.Println("Generated UUID:")
	fmt.Println(id.String())

	switch os.Args[1] {
	case "init":
		//if err := doInit(); err != nil {
		//	log.Fatalln(err)
		//}
		return
	default:
		help()
		return
	}

	//cfgFile, ok := os.LookupEnv("CONFIG")
	//if !ok {
	//    for k, v := range os.Args[1:] {
	//        if v == "-c" && len(os.Args) > k+2 {
	//            cfgFile = os.Args[k+2]
	//        }
	//        if v == "-config" && len(os.Args) > k+2 {
	//            cfgFile = os.Args[k+2]
	//        }
	//        if strings.HasPrefix(v, "-c=") {
	//            cfgFile = os.Args[k+1][3:]
	//        }
	//        if strings.HasPrefix(v, "-config=") {
	//            cfgFile = os.Args[k+1][8:]
	//        }
	//    }
	//}

	return
	//cfg := config.New()
	//flag.StringVar(&cfg.Addr, "a", "127.0.0.1:8080", "server address")
	//flag.DurationVar(&cfg.ReportInterval, "r", 10*time.Second, "interval for sending metrics to the server")
	//flag.DurationVar(&cfg.PollInterval, "p", 2*time.Second, "update metrics interval")
	//flag.StringVar(&cfg.Key, "k", "", "sign key")
	//flag.StringVar(&cfg.CryptoKey, "crypto-key", "", "path to file with public key")
	//flag.StringVar(&cfg.ConfigFile, "c", cfg.ConfigFile, "path to file with public key")
	//flag.StringVar(&cfg.ConfigFile, "config", cfg.ConfigFile, "path to file with public key")
	//flag.StringVar(&cfg.Addr, "a", cfg.Addr, "server address")
	//flag.DurationVar(&cfg.ReportInterval, "r", cfg.ReportInterval, "interval for sending metrics to the server")
	//flag.DurationVar(&cfg.PollInterval, "p", cfg.PollInterval, "update metrics interval")
	//flag.StringVar(&cfg.Key, "k", cfg.Key, "sign key")
	//flag.StringVar(&cfg.CryptoKey, "crypto-key", cfg.CryptoKey, "path to file with public key")
	//flag.Parse()

	// создадим агента по сбору и отправке метрик
	// в качестве метрик выступают различные системные характеристики машины, на которой запущен агент
	//a := agent.New(
	//	agent.WithAddress(cfg.Addr),
	//	agent.WithReportInterval(cfg.ReportInterval),
	//	agent.WithPollInterval(cfg.PollInterval),
	//	agent.WithKey(cfg.Key),
	//	agent.WithPublicKey(pubKey),
	//)

	//go http.ListenAndServe(":8091", nil) // запускаем сервер для нужд профилирования

	//a.Run()
}

func help() {
	msg := `
usage: gok [--version] [--help] <command> [<args>]

These are common GoK commands used in various situations:

start working with GoK
	init	Create a new empty database to store secret data

sync data with server
	pull 	Fetch data from server to local database
	push    Update remote data on server
`
	fmt.Println(msg)
}

func doInit() error {
	home, ok := os.LookupEnv("HOME")
	if !ok {
		msg := `
fatal: HOME environment variable not found

Please export HOME variable or use --home argument.`
		fmt.Println(msg)
		return nil
	}

	fmt.Println("home sweet home:", home)
	return nil
}

func getHome(home string) string {
	envHome, ok := os.LookupEnv("HOME")
	if ok {
		// got HOME env value, just use it
		return envHome
	}

	if !ok && home != "" {
		return home
	}

	// no proper HOME at all! sadly exit
	if !ok && home == "" {
		msg := `
fatal: HOME environment variable not found

Please export HOME variable or use --home argument.`
		fmt.Println(msg)
		os.Exit(1)
	}

	// HOME env not found, but we got home from arg
	return home
}
