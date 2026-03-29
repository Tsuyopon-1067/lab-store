package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Server struct {
		Port int `yaml:"port"`
	} `yaml:"server"`
	Database struct {
		Path string `yaml:"path"`
	} `yaml:"database"`
	Session struct {
		TimeoutMinutes int `yaml:"timeout_minutes"`
	} `yaml:"session"`
	Backup struct {
		Path     string `yaml:"path"`
		Interval int    `yaml:"interval_minutes"`
	} `yaml:"backup"`
}

func Load(configPath string) (*Config, error) {
	if configPath == "" {
		configPath = "config.yaml"
	}

	file, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	config := &Config{}
	if err := yaml.Unmarshal(file, config); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	return config, nil
}

func (c *Config) SetDefaults() {
	if c.Server.Port == 0 {
		c.Server.Port = 8080
	}
	if c.Database.Path == "" || c.Database.Path == "data/purchase.db" {
		// backend ディレクトリから実行される場合、../data/purchase.db を使用
		c.Database.Path = "../data/purchase.db"
	}
	if c.Session.TimeoutMinutes == 0 {
		c.Session.TimeoutMinutes = 30
	}
	if c.Backup.Path == "" || c.Backup.Path == "backup" {
		c.Backup.Path = "../backup"
	}
	if c.Backup.Interval == 0 {
		c.Backup.Interval = 60
	}
}
