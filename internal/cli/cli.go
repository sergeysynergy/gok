package cli

import (
	"flag"
	"fmt"
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
}

func NewCLI(logger *zap.Logger, helpMsg string) *CLI {
	return &CLI{
		lg:      logger,
		helpMsg: helpMsg,
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

	args := flag.Args()
	if len(args) == 0 {
		fmt.Println(c.helpMsg)
		return
	}

	switch args[0] {
	case "signin":
		c.signin()
	case "init":
	default:
		fmt.Println(c.helpMsg)
	}
}

func (c *CLI) signin() {
	c.lg.Debug("SIGNIN!!!")
}
