package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Server struct {
		Port string `yaml:"port"`
		Env  string `yaml:"env"`
	} `yaml:"server"`

	Database struct {
		URL          string `yaml:"url"`
		MaxOpenConns int    `yaml:"max_open_conns"`
		MaxIdleConns int    `yaml:"max_idle_conns"`
	} `yaml:"database"`

	JWT struct {
		SecretKey string `yaml:"secret_key"`
	} `yaml:"jwt"`

	Redis string `yaml:"redis"`
}

func LoadConfig(configPath string) (*Config, error) {
	config := &Config{}

	f, err := os.Open(configPath)
	if err != nil {
		return nil, err
	}

	defer f.Close()

	d := yaml.NewDecoder(f)
	if err := d.Decode(config); err != nil {
		return nil, err
	}

	return config, err
}
