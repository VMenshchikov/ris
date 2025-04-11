package app

import (
	"fmt"

	"github.com/caarlos0/env/v10"
	"github.com/ztrue/tracerr"
)

type Config struct {
	WorkerAddresses []string `env:"WORKER_ADDRESSES" envSeparator:","`
	MaxTasks        []uint   `env:"MAX_TASKS" envSeparator:","`
	WorkersId       []uint   `env:"WORKER_IDS" envSeparator:","`
}

func LoadConfig() (Config, error) {
	var cfg Config
	if err := env.Parse(&cfg); err != nil {
		return Config{}, tracerr.New(fmt.Sprintf("Ошибка загрузки конфига: %v", err))
	}

	return cfg, nil
}
