package app

import (
	"fmt"
	"log"

	"github.com/caarlos0/env"
	"github.com/ztrue/tracerr"
)

type Config struct {
	WorkerId                   uint   `env:"WORKER_ID"`
	MaxProcs                   uint   `env:"MAXPROCS"`
	KafkaBrokerUrl             string `env:"KAFKA_BROKER_URL"`
	KafkaTopicManagerToWorkers string `env:"KAFKA_TOPIC_MANAGER_TO_WORKERS"`
	KafkaTopicWorkersToManager string `env:"KAFKA_TOPIC_WORKERS_TO_MANAGER"`
}

func LoadConfig() (Config, error) {
	var cfg Config
	if err := env.Parse(&cfg); err != nil {
		return Config{}, tracerr.New(fmt.Sprintf("Ошибка загрузки конфига: %v", err))
	}
	log.Println("Config:", cfg)
	return cfg, nil
}
