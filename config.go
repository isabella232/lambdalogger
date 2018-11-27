package main

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"github.com/pkg/errors"
)

type configuration struct {
	LogConfig

	Humio       humioConfig
	TimeoutSec  int               `split_words:"true"`
	ExtraFields map[string]string `split_words:"true"`
	ExtraTags   map[string]string `split_words:"true"`
}

func (c *configuration) validate() error {
	if c.Humio.Token == "" {
		return errors.New("Must set HUMIO_TOKEN")
	}
	if c.Humio.Repository == "" {
		return errors.New("Must set HUMIO_REPOSITORY")
	}
	return nil
}

func LoadConfig(cfg interface{}) {
	if cfgFile := os.Getenv("_CONFIG_FILE"); cfgFile != "" {
		if err := godotenv.Load(cfgFile); err != nil {
			panic(errors.Wrapf(err, "Failed to load config file: %s", cfgFile))
		}
	}

	if err := envconfig.Process(os.Getenv("PREFIX"), cfg); err != nil {
		panic(errors.Wrap(err, "Failed to load initial config"))
	}
}
