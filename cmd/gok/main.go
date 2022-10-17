package main

import (
	"fmt"
	"github.com/sergeysynergy/gok/pkg/utils"
	"github.com/sergeysynergy/gok/tool/zapLogger"
	"os"
)

var (
	buildVersion string
	buildDate    string
)

func main() {
	// Check for help argument, if found or no arguments at all: display and exit.
	checkHelp()

	// Check for client version argument, if found: display and exit.
	checkVersion()

	// Init logger.
	lg := zapLogger.NewGokLogger(true)

	//var home string
	//for _, v := range os.Args {
	//	switch v {
	//	case "-":
	//	case "--home":
	//		home = "test"
	//	}
	//}

	// Arguments shift counter: +2 for every pre command argument which will be found next.
	argShift := 0
	home := getHome()
	lg.Debug("home: " + home)

	fmt.Println("SHIFT:", argShift)
	return

	switch os.Args[1] {
	case "sign-in":
		return
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
	//flag.StringVar(&cfg.Addr, "a", "127.0.0.1:8080", "service address")
	//flag.DurationVar(&cfg.ReportInterval, "r", 10*time.Second, "interval for sending metrics to the service")
	//flag.DurationVar(&cfg.PollInterval, "p", 2*time.Second, "update metrics interval")
	//flag.StringVar(&cfg.Key, "k", "", "sign key")
	//flag.StringVar(&cfg.CryptoKey, "crypto-key", "", "path to file with public key")
	//flag.StringVar(&cfg.ConfigFile, "c", cfg.ConfigFile, "path to file with public key")
	//flag.StringVar(&cfg.ConfigFile, "config", cfg.ConfigFile, "path to file with public key")
	//flag.StringVar(&cfg.Addr, "a", cfg.Addr, "service address")
	//flag.DurationVar(&cfg.ReportInterval, "r", cfg.ReportInterval, "interval for sending metrics to the service")
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

// Check if os.Args contains string argument: return -1 if not, or position in os.Args.
func argsContains(value string) int {
	for i := 1; i < len(os.Args); i++ {
		if os.Args[i] == value {
			return i
		}
	}
	return -1
}

func help() {
	msg := `
usage: gok [--help] [--version] [-h] [-u] <command> [<args>]

Pre command arguments:
	--help		List available commands.
	--version	Displays build version and date.
	-h 			Set home directory where GoK will store its files. Use HOME environment value by default.
	-u 			Set GoK user. Use USER environment value by default.

These are common GoK commands used in various situations:

start working with GoK
	init	Create new or pull existing branch to store secret data

sync data with service
	pull 	Fetch data from service to local database
	push    Update remote data on service
`
	fmt.Println(msg)
}

//func doInit() error {
//	home, ok := os.LookupEnv("HOME")
//	if !ok {
//		msg := `
//fatal: HOME environment variable not found
//
//Please export HOME variable or use --home argument.`
//		fmt.Println(msg)
//		return nil
//	}
//
//	fmt.Println("home sweet home:", home)
//	return nil
//}

func checkHelp() {
	// Checking if any argument given.
	if len(os.Args) == 1 {
		help()
		os.Exit(0)
	}

	pos := argsContains("--help")
	if pos > 0 {
		help()
		os.Exit(0)
	}
}

func checkVersion() {
	pos := argsContains("--version")
	if pos < 0 {
		return
	}

	// Add build version and date of the client binary file:
	// go run -ldflags "-X main.buildVersion=v1.0.1" main.go
	fmt.Printf("Build version: %s\n", utils.CheckNA(buildVersion))
	fmt.Printf("Build date: %s\n", utils.CheckNA(buildDate))

	os.Exit(0)
}

func getHome() string {
	var home string
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
