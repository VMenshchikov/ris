package app

import (
	"fmt"
	"log"

	"github.com/caarlos0/env"
	"github.com/ztrue/tracerr"
)

type Config struct {
	WorkerId       uint   `env:"WORKER_ID"`
	MaxProcs       uint   `env:"MAXPROCS"`
	ManagerAddress string `env:"MANAGER_ADDRESS"`
}

func LoadConfig() (Config, error) {
	var cfg Config
	if err := env.Parse(&cfg); err != nil {
		return Config{}, tracerr.New(fmt.Sprintf("Ошибка загрузки конфига: %v", err))
	}
	log.Println("Config:", cfg)
	return cfg, nil
}
