package config

import (
	"fmt"
	"io"
	"os"

	"gopkg.in/yaml.v3"
)

type Env string

type Config struct {
	Env    Env    `yaml:"env"`
	Server Server `yaml:"server"`
}

type Server struct {
	Addr string `yaml:"addr"`
	Port string `yaml:"port"`
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
	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, err
	}

	fmt.Println(*cfg)

	return cfg, nil
}
