package config

import (
	"fmt"
	"os"

	"go.yaml.in/yaml/v4"
)

type Config struct {
	App struct {
		Name    string `yaml:"name"`
		Version string `yaml:"version"`
		State   string `yaml:"state"`
	} `yaml:"app"`
	Database struct {
		Path string `yaml:"path"`
	} `yaml:"database"`
	Telegram struct {
		Token string `yaml:"token"`
	} `yaml:"telegram"`
}

func LoadConfig() (*Config, error) {
	data, err := os.ReadFile("./config/config.yaml")
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var cfg Config

	if err := yaml.Load(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	return &cfg, nil
}

func (c *Config) String() string {
	return fmt.Sprintf("App:\n"+
		"\tName: %s\n"+
		"\tVersion: %s\n"+
		"\tState: %s\n"+
		"Database:\n"+
		"\tPath: %s\n"+
		"Telegram:\n"+
		"\tToken: %s\n", c.App.Name, c.App.Version, c.App.State, c.Database.Path, c.Telegram.Token)
}
