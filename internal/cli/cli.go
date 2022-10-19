package cli

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"go.uber.org/zap"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"io/ioutil"
	"os"
	"sync"

	gokClient "github.com/sergeysynergy/gok/internal/cli/delivery/client"
	gokUC "github.com/sergeysynergy/gok/internal/cli/useCase"
	"github.com/sergeysynergy/gok/internal/data/model"
	recRepo "github.com/sergeysynergy/gok/internal/data/repository/sql/record"
	"github.com/sergeysynergy/gok/internal/entity"
	gokErrors "github.com/sergeysynergy/gok/internal/errors"
)

// Config for GoK CLI:
// AuthAddr - Auth service gRPC API address;
// StorageAddr - Storage service gRPC API address;
// Token - user auth token;
// Key - user key used to encrypt data.
type Config struct {
	mu       sync.RWMutex
	filename string

	AuthAddr    string `json:"auth_addr"`
	StorageAddr string `json:"storage_addr"`
	Token       string `json:"token"`
	Key         string `json:"key"`
	Branch      string `json:"branch"`
	LocalHead   uint64 `json:"head"`
}

func NewConf(filename string) *Config {
	const (
		defaultAuthAddr    = ":7000"
		defaultStorageAddr = ":7001"
	)
	cfg := &Config{
		AuthAddr:    defaultAuthAddr,
		StorageAddr: defaultStorageAddr,
		filename:    filename,
	}

	// Trying to read config file.
	if err := cfg.Read(); err != nil {
		fmt.Println("Failed to read config file -", err.Error())
	}

	return cfg
}

// Write config struct to json file.
func (c *Config) Write() error {
	// TODO: rewrite config saving service addresses.
	c.mu.Lock()
	defer c.mu.Unlock()

	jsonString, _ := json.Marshal(c)
	err := ioutil.WriteFile(c.filename, jsonString, os.ModePerm)
	if err != nil {
		return err
	}

	return nil
}

// Read config struct from json file.
func (c *Config) Read() error {
	c.mu.RLock()
	defer c.mu.RUnlock()

	data, err := os.ReadFile(c.filename)
	if err != nil {
		return fmt.Errorf("error when opening file: %w", err)
	}

	var cfg Config
	err = json.Unmarshal(data, &cfg)
	if err != nil {
		return fmt.Errorf("error during Unmarshal(): %w", err)
	}

	c.AuthAddr = cfg.AuthAddr
	c.StorageAddr = cfg.StorageAddr
	c.Token = cfg.Token
	c.Key = cfg.Key
	c.Branch = cfg.Branch
	c.LocalHead = cfg.LocalHead

	return nil
}

// CLI contains argument values and methods for command line processing.
type CLI struct {
	lg *zap.Logger
	// Help message
	helpMsg string
	// Config for CLI.
	cfg *Config
	// Directory to store GoK local files.
	home string
	// Username
	user string
	// All CLI arguments goes after flags arguments.
	args []string
	// gRPC client to access GoK API.
	client *gokClient.Client
	// Use cases to work with GoK API.
	uc gokUC.UseCase
	// Database creds.
	dbOnce *sync.Once
	db     *gorm.DB
}

func New(logger *zap.Logger, helpMsg string) *CLI {
	cli := &CLI{
		lg:      logger,
		helpMsg: helpMsg,
		dbOnce:  &sync.Once{},
	}
	cli.initCLI()

	return cli
}

func (c *CLI) initCLI() {
	if err := c.parsePreCommandsArgs(); err != nil {
		c.lg.Fatal(err.Error())
	}

	// Create user directory.
	dir := c.home + "/.gok"
	if err := os.Mkdir(c.home+"/.gok", os.ModePerm); err != nil {
		if err.Error() != "mkdir "+dir+": file exists" {
			c.lg.Fatal(err.Error())
			return
		}
	}
	dir = c.home + "/.gok/" + c.user
	if err := os.Mkdir(dir, os.ModePerm); err != nil {
		if err.Error() != "mkdir "+dir+": file exists" {
			c.lg.Fatal(err.Error())
			return
		}
	}

	// Init GoK user config or read existing one from file.
	filename := c.home + "/.gok/config.json"
	c.cfg = NewConf(filename)

	c.dbConnect()
	c.newUseCase()
}

// parsePreCommandsArgs extract home and username info and create config file [home]/.gok/config.json
func (c *CLI) parsePreCommandsArgs() error {
	c.home, _ = os.LookupEnv("HOME")
	c.user, _ = os.LookupEnv("USER")

	flag.StringVar(&c.home, "h", c.home, "set home directory where GoK will store its files")
	flag.StringVar(&c.user, "u", c.user, "set GoK user")
	flag.Parse()

	if c.home == "" {
		msg := `Failed using cli: need to define home directory where cli can store local files.
It could be HOME environment value. Or you can redefine it using -h flag.`
		fmt.Println(msg)
		return fmt.Errorf("home not found")
	}

	if c.user == "" {
		msg := `Failed using cli: need to define cli username.
It could be USER environment value. Or you can redefine it using -u flag.`
		fmt.Println(msg)
		return fmt.Errorf("username not found")
	}

	c.lg.Debug("got HOME: " + c.home)
	c.lg.Debug("got USER: " + c.user)
	return nil
}

func (c *CLI) dbConnect() {
	c.dbOnce.Do(func() {
		dbPath := fmt.Sprintf("%s/.gok/%s/default.db", c.home, c.user)
		db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
		if err != nil {
			c.lg.Fatal(fmt.Sprintf("connection to SQLite failed: %s", err))
		}

		// Create and migrate database tables.
		err = db.AutoMigrate(&model.Record{})
		if err != nil {
			c.lg.Fatal(fmt.Sprintf("auto migration has failed: %s", err))
		}

		c.db = db
		c.lg.Info("established connection with DB")
	})
}

func (c *CLI) newUseCase() {
	client := gokClient.New(c.lg, c.cfg.AuthAddr, c.cfg.StorageAddr)
	repo := recRepo.New(c.db)
	c.uc = gokUC.New(c.lg, repo, client)
}

// Parse method to process main CLI commands: something like router.
func (c *CLI) Parse() {
	c.args = flag.Args()
	if len(c.args) == 0 {
		fmt.Println(c.helpMsg)
		return
	}

	switch c.args[0] {
	case "signin":
		c.signIn()
	case "login":
		c.login()
	case "init":
		c.init()
	case "desc":
		c.desc()
	case "push":
		c.push()
	default:
		fmt.Println(c.helpMsg)
	}
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
///   Commands section   ///////////////////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

// Method signIn create new user record locally and at the server side.
func (c *CLI) signIn() {
	if len(c.args) > 1 {
		fmt.Println("Invalid argument: home and user enough to execute signin.")
		return
	}

	usr := &entity.CLIUser{
		Login: c.user,
		Home:  c.home,
	}

	signedUsr, err := c.uc.SignIn(usr)
	if err != nil {
		if errors.Is(err, gokErrors.ErrUserAlreadyExists) {
			fmt.Println("User already exists. Try to login.")
		} else {
			c.lg.Error(err.Error())
		}
		return
	}

	c.cfg.Token = signedUsr.Token
	c.cfg.Key = signedUsr.Key
	if err = c.cfg.Write(); err != nil {
		c.lg.Error(err.Error())
		return
	}

	fmt.Println("New user has been successfully registered. Now init new branch to store your secrets.")
}

// Method login get user token for auth process.
func (c *CLI) login() {
	if len(c.args) > 1 {
		fmt.Println("Invalid arguments for login.")
		return
	}

	usr := &entity.CLIUser{
		Login: c.user,
		Home:  c.home,
	}

	signedUsr, err := c.uc.Login(usr)
	if err != nil {
		c.lg.Error(err.Error())
		if errors.Is(err, gokErrors.ErrUserNotFound) {
			fmt.Println("User not found: signin first.")
			return
		}
		fmt.Println("Login failed, try to signin first.")
		return
	}

	c.cfg.Token = signedUsr.Token
	c.cfg.Key = signedUsr.Key
	if err = c.cfg.Write(); err != nil {
		c.lg.Error(err.Error())
		return
	}

	fmt.Println("Login successful. Now init new branch to store your secrets.")
}

// Method init add or get branch info and store it locally in config file.
func (c *CLI) init() {
	if len(c.args) > 1 {
		fmt.Println("Invalid argument: home and user enough to execute init... so far.")
		return
	}

	brn, err := c.uc.Init(c.cfg.Token)
	if err != nil {
		if errors.Is(err, gokErrors.ErrAuthRequired) {
			fmt.Println("Authentication required: try to signin or login.")
		} else {
			c.lg.Error(err.Error())
		}
		return
	}

	c.cfg.Branch = brn.Name

	if brn.Head > c.cfg.LocalHead {
		c.cfg.LocalHead = brn.Head
		// TODO: add git pull command
	}

	if err = c.cfg.Write(); err != nil {
		c.lg.Error(err.Error())
		return
	}

	fmt.Println("Branch `" + brn.Name + "` has been successfully initiated. Now it's time to add some secrets!")
}

func (c *CLI) push() {
	if len(c.args) > 1 {
		fmt.Println("Too many arguments for push.")
		return
	}

	brn, err := c.uc.Push(c.cfg.Token, c.cfg.Branch, c.cfg.LocalHead)
	if err != nil {
		if errors.Is(err, gokErrors.ErrLocalBranchBehind) {
			fmt.Println("Your local branch is behind server: please make pull first to update data.")
			return
		}
		if errors.Is(err, gokErrors.ErrRecordNotFound) {
			fmt.Println("No new records for push.")
			return
		}
		if errors.Is(err, gokErrors.ErrAuthRequired) {
			fmt.Println("Authentication required: try to signin or login.")
			return
		}

		c.lg.Error(err.Error())
		return
	}
	if brn == nil {
		err = gokErrors.ErrPushUnknown
		c.lg.Error(err.Error())
		fmt.Println("Failed to push -", err)
		return
	}

	if brn.Head > c.cfg.LocalHead && brn.Name == c.cfg.Branch {
		c.cfg.LocalHead = brn.Head
		if err = c.cfg.Write(); err != nil {
			c.lg.Error(err.Error())
			return
		}
	}

	c.lg.Debug(fmt.Sprintf("local branch header: %d", c.cfg.LocalHead))
	fmt.Println("Push successful.")
}
