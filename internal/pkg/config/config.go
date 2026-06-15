package config

import (
	"os"
	"strconv"

	"github.com/davidsugianto/go-pkgs/config"
)

var (
	cfg *Config
)

func Load(path string) (*Config, error) {
	if path == "" {
		path = "configs/config.yaml"
	}
	cfg, err := config.Load[Config](path)
	if err != nil {
		return nil, err
	}

	cfg.applyEnvOverrides()

	return &cfg, nil
}

func GetConfig() *Config {
	if cfg == nil {
		cfg = &Config{}
	}
	return cfg
}

func (c *Config) applyEnvOverrides() {
	if v := os.Getenv("PORT"); v != "" {
		if port, err := strconv.Atoi(v); err == nil {
			c.Server.Port = port
		}
	}
}

type Config struct {
	Server ServerConfig `json:"server" yaml:"server"`
	CORS   CORSConfig   `json:"cors" yaml:"cors"`
}

type ServerConfig struct {
	Port int `json:"port" yaml:"port"`
}

type CORSConfig struct {
	AllowedOrigins   []string `json:"allowed_origins" yaml:"allowed_origins"`
	AllowedMethods   []string `json:"allowed_methods" yaml:"allowed_methods"`
	AllowedHeaders   []string `json:"allowed_headers" yaml:"allowed_headers"`
	AllowCredentials bool     `json:"allow_credentials" yaml:"allow_credentials"`
}
