package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type (
	Config struct {
		Server       ServerConfig `yaml:"server"`
		Database     DBConfig     `yaml:"database"`
		SMTP         SMTPConfig   `yaml:"smtp"`
		AWS          AWSConfig    `yaml:"aws"`
		Redis        string       `yaml:"redis"`
		JWTSecretKey string       `yaml:"jwt_secret_key"`
	}

	ServerConfig struct {
		Port string `yaml:"port"`
		Env  string `yaml:"env"`
	}

	DBConfig struct {
		URL          string `yaml:"url"`
		MaxOpenConns int    `yaml:"max_open_conns"`
		MaxIdleConns int    `yaml:"max_idle_conns"`
	}

	SMTPConfig struct {
		Host string `yaml:"host"`
		Port string `yaml:"port"`
		From string `yaml:"from"`
	}

	AWSConfig struct {
		Region string `yaml:"region"`
		Bucket string `yaml:"bucket"`
	}
)

func LoadConfig(configPath string) (*Config, error) {
	config := &Config{}

	root, err := os.OpenRoot(".")
	if err != nil {
		return nil, fmt.Errorf("failed to open config: %w", err)
	}

	f, err := root.Open(configPath)
	if err != nil {
		return nil, err
	}

	defer func() {
		if cerr := f.Close(); cerr != nil {
			err = fmt.Errorf("failed to close file: %w", err)
		}
	}()

	d := yaml.NewDecoder(f)
	if err := d.Decode(config); err != nil {
		return nil, err
	}

	config.Server.Port = getEnv("SERVER_PORT", config.Server.Port)
	config.Server.Env = getEnv("APP_ENV", config.Server.Env)
	config.Database.URL = getEnv("DB_URL", config.Database.URL)
	config.JWTSecretKey = getEnv("JWT_SECRET", config.JWTSecretKey)
	config.Redis = getEnv("REDIS_URL", config.Redis)

	return config, err
}

func getEnv(key, currentVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return currentVal
}
