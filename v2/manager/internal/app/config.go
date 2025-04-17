package app

import (
	"fmt"

	"github.com/caarlos0/env/v10"
	"github.com/ztrue/tracerr"
)

type Environment struct {
	KafkaBrokerUrl             string `env:"KAFKA_BROKER_URL"`
	KafkaTopicManagerToWorkers string `env:"KAFKA_TOPIC_MANAGER_TO_WORKERS"`
	KafkaTopicWorkersToManager string `env:"KAFKA_TOPIC_WORKERS_TO_MANAGER"`
	MongoUri                   string `env:"MONGO_URI"`
	IsDebug                    bool   `env:"IS_DEBUG"`
}
type Config struct {
	Env Environment
}

func LoadConfig() (*Config, error) {
	var cfg Environment
	if err := env.Parse(&cfg); err != nil {
		return nil, tracerr.New(fmt.Sprintf("Ошибка загрузки конфига: %v", err))
	}

	return &Config{
		Env: cfg,
	}, nil
}
