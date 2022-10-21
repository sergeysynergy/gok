package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/sergeysynergy/gok/internal/cli"
	"github.com/sergeysynergy/gok/pkg/utils"
	"github.com/sergeysynergy/gok/tool/zapLogger"
)

var (
	buildVersion string
	buildDate    string
	debug        string
)

func main() {
	// Init logger.
	if debug == "" {
		debug = "true"
	}
	dbg, err := strconv.ParseBool(debug)
	if err != nil {
		log.Fatalln("Failed to convert debug value -", err)
	}
	lg := zapLogger.NewGokLogger(dbg)

	// Check for help argument, if found or no arguments at all: display and exit.
	var helpMsg string
	err = checkHelp(&helpMsg)
	if err != nil {
		lg.Debug(err.Error())
		fmt.Println(helpMsg)
		return
	}

	// Check for client version argument, if found: display and exit.
	err = checkVersion()
	if err != nil {
		lg.Debug(err.Error())
		return
	}

	// Define structure with values and methods for arguments processing needed for GoK.
	cli := cli.New(lg, helpMsg)

	// Perform the main amount of argument parsing and further operations.
	cli.Parse()
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

func checkHelp(msg *string) error {
	*msg = `
usage: cli [--help] [--version] [-h] [-u] <command> [<args>]

Pre command arguments:
	--help		List available commands.
	--version	Displays build version and date.
	-hm			Set home directory where GoK will store its files. Use HOME environment value by default.
	-u 			Set GoK user. Use USER environment value by default.

These are common GoK commands used in various situations:

start working with GoK
	init	Create new or pull existing branch to store secret data

	desc add		Create new description record
	desc ls			List all description records
	desc set [ID]	Update description record by given ID

	sync data with service
		pull 	Fetch data from service to local database
		push    Update remote data on service`

	// Checking if any argument given.
	if len(os.Args) == 1 {
		return fmt.Errorf("no arguments given")
	}

	pos := argsContains("--help")
	if pos > 0 {
		return fmt.Errorf("found help request")
	}

	return nil
}

func checkVersion() error {
	pos := argsContains("--version")
	if pos < 0 {
		return nil
	}

	// Add build version and date of the client binary file:
	// go run -ldflags "-X main.buildVersion=v1.0.1" main.go
	fmt.Printf("Build version: %s\n", utils.CheckNA(buildVersion))
	fmt.Printf("Build date: %s\n", utils.CheckNA(buildDate))

	return fmt.Errorf("found version request")
}
