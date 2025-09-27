package database

type Config struct {
	Driver string `yaml:"driver"`
	Addr   string `yaml:"addr"`
	Port   string `yaml:"port"`
	DB     string `yaml:"db"`
	User   string `yaml:"user"   envconfig:"DB_USER"`
	Pass   string `yaml:"pass"   envconfig:"DB_PASS"`
}
