package entity

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"sync"
)

// CLIConf for GoK CLI:
// AuthAddr - Auth service gRPC API address;
// StorageAddr - Storage service gRPC API address;
// Token - user auth token;
// Key - user key used to encrypt data.
type CLIConf struct {
	mu       sync.RWMutex
	filename string

	AuthAddr    string `json:"auth_addr"`
	StorageAddr string `json:"storage_addr"`
	Token       string `json:"token"`
	Key         string `json:"key"`
	BranchID    uint64 `json:"branch_id"`
	LocalHead   uint64 `json:"head"`
}

func NewCLIConf(filename string) *CLIConf {
	const (
		defaultAuthAddr    = ":7000"
		defaultStorageAddr = ":7001"
	)
	cfg := &CLIConf{
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
func (c *CLIConf) Write() error {
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
func (c *CLIConf) Read() error {
	c.mu.RLock()
	defer c.mu.RUnlock()

	data, err := os.ReadFile(c.filename)
	if err != nil {
		return fmt.Errorf("error when opening file: %w", err)
	}

	var cfg CLIConf
	err = json.Unmarshal(data, &cfg)
	if err != nil {
		return fmt.Errorf("error during Unmarshal(): %w", err)
	}

	c.AuthAddr = cfg.AuthAddr
	c.StorageAddr = cfg.StorageAddr
	c.Token = cfg.Token
	c.Key = cfg.Key
	c.BranchID = cfg.BranchID
	c.LocalHead = cfg.LocalHead

	return nil
}
