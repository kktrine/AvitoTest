package config

import (
	"flag"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	//LogPath  string   `yaml:"log_path" env-required:"true"`
	DbConfig DbConfig `yaml:"db" env-required:"true"`
}

type DbConfig struct {
	Host     string `yaml:"host" env-default:"localhost"`
	Port     int    `yaml:"port" env-default:"5432"`
	User     string `yaml:"user" env-required:"true"`
	Password string `yaml:"password" env-required:"true"`
	DbName   string `yaml:"db" env-required:"true"`
	SSL      string `yaml:"ssl" env-default:"disable"`
}

func MustLoad() *Config {
	path := getCfgPath()
	if _, err := os.Stat(path); os.IsNotExist(err) {
		panic("config file doesn't exist" + err.Error())
	}
	var cfg Config
	if err := cleanenv.ReadConfig(path, &cfg); err != nil {
		panic("fail to read config file" + err.Error())
	}
	return &cfg
}

func getCfgPath() string {
	var path string
	flag.StringVar(&path, "config", "./config/config.yaml", "path to config file")
	flag.Parse()
	return path
}
