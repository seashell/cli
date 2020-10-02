package dev

import (
	"os"
	"os/user"
	"path"

	"github.com/hashicorp/hcl/v2/hclsimple"
	"github.com/seashell/cli/version"
)

// Config contains configurations for the Seashell CLI
type Config struct {
	// LogLevel is the level of the logs to put out
	LogLevel string `hcl:"log_level,optional"`

	// DataDir is the directory to store our state in
	DataDir string

	// ProjectID
	ProjectID string

	// ProjectDir
	ProjectDir string

	// VagrantDir = DataDir + ProjectID
	VagrantDir string

	// Version information (set at compilation time)
	Version *version.VersionInfo
}

// Merge merges two Config structs, returning the result
func (c *Config) Merge(b *Config) *Config {

	if b == nil {
		return c
	}

	result := *c

	if b.LogLevel != "" {
		result.LogLevel = b.LogLevel
	}
	if b.DataDir != "" {
		result.DataDir = b.DataDir
	}
	if b.ProjectID != "" {
		result.ProjectID = b.ProjectID
	}

	if b.ProjectDir != "" {
		result.ProjectDir = b.ProjectDir
	}
	if b.Version != nil {
		result.Version = b.Version
	}

	return &result
}

// DefaultConfig returns a Config struct populated with sane defaults
func DefaultConfig() *Config {

	usr, err := user.Current()
	if err != nil {
		panic(err)
	}

	CurrentDir, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	defaultDataDir := path.Join(usr.HomeDir, ".seashell")

	return &Config{
		LogLevel:   "INFO",
		DataDir:    defaultDataDir,
		ProjectID:  "",
		ProjectDir: CurrentDir,
		VagrantDir: "",
		Version:    version.GetVersion(),
	}
}

// Validate returns an error in case a Config struct is invalid.
func (c *Config) Validate() error {
	// TODO: implement validation
	return nil
}

// LoadFromFile loads the configuration from a given path
func (c *Config) LoadFromFile(path string) (*Config, error) {

	_, err := os.Stat(path)
	if err != nil {
		return nil, err
	}

	config := &Config{}

	err = hclsimple.DecodeFile(path, nil, config)
	if err != nil {
		return nil, err
	}

	return config, nil
}
