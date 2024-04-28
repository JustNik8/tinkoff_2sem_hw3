package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type ServerConfig struct {
	Postgres PostgresConfig   `yaml:"postgres"`
	Server   ServerInfoConfig `yaml:"server"`
	Kafka    KafkaConfig      `yaml:"kafka"`
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

type KafkaConfig struct {
	Addrs []string `yaml:"addrs"`
	Topic string   `yaml:"topic"`
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
