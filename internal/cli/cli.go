package cli

import (
	"flag"
	"fmt"
	gokClient "github.com/sergeysynergy/gok/internal/cli/client"
	"go.uber.org/zap"
	"os"
)

// CLI contains argument values and methods for command line processing.
type CLI struct {
	lg *zap.Logger
	// Help message
	helpMsg string
	// Directory to store GoK local files.
	home string
	// Username
	user string
	// All CLI arguments goes after flags arguments.
	args []string
	// gRPC client to access GoK API.
	client *gokClient.Client
}

func New(logger *zap.Logger, helpMsg, authAddr, storageAddr string) *CLI {
	return &CLI{
		lg:      logger,
		helpMsg: helpMsg,
		client:  gokClient.New(logger, authAddr, storageAddr),
	}
}

func (c *CLI) preCommandsCheck() {
	c.home, _ = os.LookupEnv("HOME")
	c.user, _ = os.LookupEnv("USER")

	flag.StringVar(&c.home, "h", c.home, "set home directory where GoK will store its files")
	flag.StringVar(&c.user, "u", c.user, "set GoK user")
	flag.Parse()

	if c.home == "" {
		msg := `Failed using cli: need to define home directory where cli can store local files.
It could be HOME environment value. Or you can redefine it using -h flag.`
		fmt.Println(msg)
		os.Exit(0)
	}

	if c.user == "" {
		msg := `Failed using cli: need to define cli username.
It could be USER environment value. Or you can redefine it using -u flag.`
		fmt.Println(msg)
		os.Exit(0)
	}

	c.lg.Debug("got HOME: " + c.home)
	c.lg.Debug("got USER: " + c.user)
}

func (c *CLI) Parse() {
	c.preCommandsCheck()

	c.args = flag.Args()
	if len(c.args) == 0 {
		fmt.Println(c.helpMsg)
		return
	}

	switch c.args[0] {
	case "signin":
		c.signIn()
	case "init":
	default:
		fmt.Println(c.helpMsg)
	}
}

func (c *CLI) signIn() {
	if len(c.args) > 1 {
		fmt.Println("Invalid argument: home and user enough to execute signin.")
		os.Exit(0)
	}
	c.lg.Debug("SIGNIN!!!")
}
