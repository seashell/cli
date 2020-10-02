package command

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/caarlos0/env"
	"github.com/joho/godotenv"
	"github.com/seashell/cli/dev"
	cli "github.com/seashell/cli/pkg/cli"
	"github.com/seashell/cli/pkg/log"
	"github.com/seashell/cli/pkg/log/zap"
)

// DevListCommand :
type DevListCommand struct {
	UI cli.UI
}

// Name :
func (c *DevListCommand) Name() string {
	return "dev list"
}

// Synopsis :
func (c *DevListCommand) Synopsis() string {
	return "List configured Seashell development environments"
}

// Run :
func (c *DevListCommand) Run(ctx context.Context, args []string) int {

	config := c.parseConfig(args)
	config = dev.DefaultConfig().Merge(config)

	logger, err := zap.NewLoggerAdapter(zap.Config{
		LoggerOptions: log.LoggerOptions{
			Level:  config.LogLevel,
			Prefix: "env: ",
		},
	})

	if err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	d, err := dev.New(config, logger)
	if err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	err = d.List()
	if err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	return 0
}

func (c *DevListCommand) parseConfig(args []string) *dev.Config {

	flags := FlagSet(c.Name())

	configFromFlags := c.parseFlags(flags, args)
	configFromFile := c.parseConfigFiles(flags.configPaths...)
	configFromEnv := c.parseEnv(flags.envPaths...)

	config := &dev.Config{}

	config = config.Merge(configFromFile)
	config = config.Merge(configFromEnv)
	config = config.Merge(configFromFlags)

	if err := config.Validate(); err != nil {
		c.UI.Error(fmt.Sprintf("Invalid input: %s", err.Error()))
		os.Exit(1)
	}

	return config
}

func (c *DevListCommand) parseFlags(flags *RootFlagSet, args []string) *dev.Config {

	flags.Usage = func() {
		c.UI.Output("\n" + c.Help() + "\n")
	}

	config := &dev.Config{}

	// General options
	flags.StringVar(&config.DataDir, "data-dir", "", "")
	flags.StringVar(&config.LogLevel, "log-level", "", "")

	if err := flags.Parse(args); err != nil {
		c.UI.Error("==> Error: " + err.Error() + "\n")
		os.Exit(1)
	}

	return config
}

func (c *DevListCommand) parseConfigFiles(paths ...string) *dev.Config {

	config := &dev.Config{}

	if len(paths) > 0 {
		// TODO : Load configurations from HCL files
		c.UI.Info(fmt.Sprintf("==> Loading configurations from: %v", paths))
	}

	return config
}

func (c *DevListCommand) parseEnv(paths ...string) *dev.Config {

	config := &dev.Config{}

	if len(paths) > 0 {

		c.UI.Info(fmt.Sprintf("==> Loading environment variables from: %v", paths))
		c.UI.Warn(fmt.Sprintf("  - This will not override already existing variables!"))

		err := godotenv.Load(paths...)

		if err != nil {
			c.UI.Error(fmt.Sprintf("Error parsing env files: %s", err.Error()))
			os.Exit(1)
		}
	}

	env.Parse(config)

	return config
}

// Help :
func (c *DevListCommand) Help() string {
	h := `
Usage: seashell dev destroy [options]
	
  Lists all current configured Seashell development environments.

General Options:
` + GlobalOptions() + `

List Options:

  --log-level=<level>
    The logging level Seashell CLI should log at. Valid values are INFO, WARN, DEBUG, ERROR, FATAL.	
`
	return strings.TrimSpace(h)
}
