package database

type Config struct {
	Driver string `yaml:"driver"`
	DSN    string `yaml:"dsn"`
}
