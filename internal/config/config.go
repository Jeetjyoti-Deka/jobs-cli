package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Company struct {
	Name     string `yaml:"name"`
	Provider string `yaml:"provider"`
	Board    string `yaml:"board"`
}

type Config struct {
	Companies []Company `yaml:"companies"`
}

func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg Config

	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}