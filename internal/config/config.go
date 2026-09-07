package config

import (
	"os"

	"github.com/zro404/lattice/internal/logger"
	"gopkg.in/yaml.v3"
)

type YamlConfig struct {
	Version int `yaml:"version"`

	Databases []struct {
		Name string `yaml:"name"`
	} `yaml:"databases"`

	Logs []struct {
		Name  string   `yaml:"name"`
		Paths []string `yaml:"paths"`
	} `yaml:"logs"`
}

func LoadFile(path string) *YamlConfig {
	stream, err := os.ReadFile(path)
	if err != nil {
		logger.Fatalf("Error reading config file: %s", err.Error())
	}

	var config YamlConfig

	if err := yaml.Unmarshal(stream, &config); err != nil {
		logger.Fatalf("Config Error: %s", err.Error())
	}

	return &config
}
