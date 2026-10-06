package config

import (
	"bytes"
	"errors"
	"os"
	"time"

	"github.com/zro404/lattice/internal/logger"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Version int `yaml:"version"`

	TickInterval time.Duration `yaml:"tick_interval"`

	Databases []Database `yaml:"databases"`

	Logs []Log `yaml:"logs"`
}

type Database struct {
	Name string `yaml:"name"`
}

type Log struct {
	Name  string   `yaml:"name"`
	Paths []string `yaml:"paths"`
}

func configError(msg string) error {
	return errors.New("Config Error: " + msg)
}

func (c *Config) Validate() error {
	if c.Version <= 0 {
		return configError("Invalid config version")
	}

	if c.TickInterval < time.Second {
		return configError("Invalid tick_interval (must be >= 1s) [format: 1h5m30s]")
	}

	for _, db := range c.Databases {
		if db.Name == "" {
			return configError("Invalid database name")
		}
	}

	for _, log := range c.Logs {
		if log.Name == "" {
			return configError("Invalid log name")
		}

		if len(log.Paths) == 0 {
			return configError("Log paths cannot be empty")
		}
	}

	return nil
}

func LoadFile(path string) (*Config, error) {
	stream, err := os.ReadFile(path)
	if err != nil {
		return nil, configError("Error reading config file: " + err.Error())
	}

	var config Config

	decoder := yaml.NewDecoder(bytes.NewReader(stream))
	decoder.KnownFields(true)

	if err := decoder.Decode(&config); err != nil {
		return nil, configError(err.Error())
	}

	if err := config.Validate(); err != nil {
		logger.Fatalf("%s", err.Error())
	}

	return &config, nil
}
