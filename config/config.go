package config

import (
	"io"
	"os"
	"time"

	"github.com/kelseyhightower/envconfig"
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
	DB      DB      `yaml:"db"`
	ML      ML      `yaml:"ml"`
}

type Server struct {
	Addr         string        `yaml:"addr" envconfig:"SERVER_ADDR"`
	Port         string        `yaml:"port" envconfig:"SERVER_PORT"`
	ReadTimeout  time.Duration `yaml:"read_timeout"`
	WriteTimeout time.Duration `yaml:"write_timeout"`
}

type ML struct {
	Addr string `yaml:"addr" envconfig:"ML_ADDR"`
	Port string `yaml:"port" envconfig:"ML_PORT"`
}

type Sensors struct {
	HandshakeTimeout time.Duration  `yaml:"handshake_timeout"`
	Entities         []SensorEntity `yaml:"entities"`
}

type SensorEntity struct {
	UUID  string `yaml:"uuid"`
	Token string `yaml:"token"`
	IP    string `yaml:"ip"`
}

type DB struct {
	Driver string `yaml:"driver"`
	DSN    string `yaml:"dsn"`
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

	// Читаем переменные окружения
	if err := envconfig.Process("", cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}
