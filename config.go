package main

import (
	"errors"
)

type configuration struct {
	LogConfig

	Humio      humioConfig
	TimeoutSec int `split_words:"true"`
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
