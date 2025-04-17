package app

import (
	"hash_worker/internal/client"
	consumer "hash_worker/internal/interface/kafka"
	"hash_worker/internal/services/alive"
	"hash_worker/internal/usecases"
	"log"
	"runtime"

	"github.com/segmentio/kafka-go"
)

type App struct {
	config       Config
	consumer     *consumer.Consumer
	sender       *client.Sender
	aliveService *alive.AliveService
}

func NewApp(cfg Config) *App {
	writer, reader, err := connectKafka(cfg.KafkaBrokerUrl, cfg.KafkaTopicManagerToWorkers, cfg.KafkaTopicWorkersToManager)
	if err != nil {
		log.Panicln(err)
	}

	sender := client.New(writer)
	log.Println(writer, sender)
	crack := usecases.Crack{
		WorkerID: cfg.WorkerId,
		Sender:   sender,
	}
	consumer := consumer.New(reader, &crack)

	app := App{
		config:   cfg,
		consumer: consumer,
		sender:   sender,
		aliveService: &alive.AliveService{
			Sender:   sender,
			WorkerId: cfg.WorkerId,
			MaxTasks: cfg.MaxProcs,
		},
	}

	app.setProcs()

	return &app
}

func connectKafka(brokerURL, managerToWorkersTopic, workersToManagerTopic string) (*kafka.Writer, *kafka.Reader, error) {
	// Настройка Kafka Writer
	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers:      []string{brokerURL},   // Адрес Kafka-брокера
		Topic:        workersToManagerTopic, // Тема для отправки сообщений от менеджера к воркерам
		Balancer:     &kafka.LeastBytes{},   // Балансировщик для выбора партиции
		RequiredAcks: int(kafka.RequireAll),
		BatchSize:    1,
	})

	// Настройка Kafka Reader
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:        []string{brokerURL},   // Адрес Kafka-брокера
		Topic:          managerToWorkersTopic, // Тема для получения сообщений от воркеров
		GroupID:        "worker-group",        // ID группы потребителей
		CommitInterval: 0,
		StartOffset:    kafka.LastOffset,
	})

	log.Println("Connected to Kafka successfully")
	return writer, reader, nil
}

func (a *App) setProcs() {
	if a.config.MaxProcs < uint(runtime.NumCPU()) {
		runtime.GOMAXPROCS(int(a.config.MaxProcs))
	}
	log.Printf("Procs: %d/%d", a.config.MaxProcs, runtime.NumCPU())
}

func (a *App) RunApp() {
	log.Println("Сервер запущен")
	go a.aliveService.TranslateAlive()
	a.consumer.Run()
}
