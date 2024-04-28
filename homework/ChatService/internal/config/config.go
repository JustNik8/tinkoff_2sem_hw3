package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type ServerConfig struct {
	Postgres PostgresConfig   `yaml:"postgres"`
	Server   ServerInfoConfig `yaml:"server"`
	Redis    RedisConfig      `yaml:"redis"`
}

type PostgresConfig struct {
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Database string `yaml:"database"`
}

type ServerInfoConfig struct {
	Port int `yaml:"port"`
}

type ClientConfig struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

type RedisConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Password string `yaml:"password"`
	DB       int    `yaml:"db"`
}

func ParseServerConfig(path string) (*ServerConfig, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	cfg := &ServerConfig{}
	err = yaml.NewDecoder(f).Decode(cfg)

	return cfg, err
}

func ParseClientConfig(path string) (*ClientConfig, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	cfg := &ClientConfig{}
	err = yaml.NewDecoder(f).Decode(cfg)

	return cfg, err
}
