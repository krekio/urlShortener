package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"os"
	"time"
)

type Config struct {
	Env        string `env:"ENV" envDefault:"local"`
	StorageDsn string `env:"STORAGE_DSN" env-required:"true"`
	ConfigPath string `env:"CONFIG_PATH" envDefault:"./.env" env-required:"true"`
}

type HTTPServer struct {
	Address     string        `env:"HTTP_ADDR" env-default:"localhost:8080"`
	Timeout     time.Duration `env:"HTTP_TIMEOUT" env-default:"5s"`
	IdleTimeout time.Duration `env:"HTTP_IDLE_TIME" env-default:"60s"`
	User        string        `env:"HTTP_USER" env-required:"true"`
	Password    string        `env:"HTTP_PASSWORD" env-required:"true"`
}

func NewConfig() (*Config, *HTTPServer) {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		log.Fatal("CONFIG_PATH environment variable not set")
	}
	if _, err := os.Stat(configPath); err != nil {
		log.Fatal("CONFIG_PATH does not exist")
	}
	var cfg Config
	var httpServer HTTPServer
	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatal(err)
	}
	if err := cleanenv.ReadConfig(configPath, &httpServer); err != nil {
		log.Fatal(err)
	}
	return &cfg, &httpServer
}
