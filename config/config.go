package config

import (
	"io"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

type Env string

func (e Env) String() string { return string(e) }

const (
	EnvDev Env = "dev"
)

type Config struct {
	Env     Env     `yaml:"env"`
	Server  Server  `yaml:"server"`
	Sensors Sensors `yaml:"sensors"`
}

type Server struct {
	Addr         string        `yaml:"addr"`
	Port         string        `yaml:"port"`
	ReadTimeout  time.Duration `yaml:"read_timeout"`
	WriteTimeout time.Duration `yaml:"write_timeout"`
}

type Sensors struct {
	HandshakeTimeout time.Duration  `yaml:"handshake_timeout"`
	Entities         []SensorEntity `yaml:"entities"`
}

type SensorEntity struct {
	UUID  string `yaml:"uuid"`
	Token string `yaml:"token"`
}

func ReadConfig(path string) (*Config, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	cfg := &Config{}
	if err = yaml.Unmarshal(data, cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}
