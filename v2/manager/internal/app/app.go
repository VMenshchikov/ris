package app

import (
	"context"
	"fmt"
	"hash_manager/internal/infra/storage"
	"hash_manager/internal/interface/httpserver"
	kafka_consumer "hash_manager/internal/interface/kafka"
	sch "hash_manager/internal/services/scheduler"
	monitor "hash_manager/internal/services/workers_monitor"
	"hash_manager/internal/usecases"
	"log"
	"time"

	myKafka "hash_manager/internal/infra/kafka"

	"github.com/segmentio/kafka-go"
	"github.com/ztrue/tracerr"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type App struct {
	config        *Config
	server        *httpserver.Server
	scheduler     *sch.Scheduler
	kafkaConsumer *kafka_consumer.KafkaConsumer
	kafkaSender   *myKafka.Sender
	//kafkaReader *kafka.Reader
	//kafkaWriter *kafka.Writer
}

func NewApp(cfg *Config) (*App, error) {

	// Подключение к Kafka
	kafkaWriter, kafkaReader, err := connectKafka(cfg.Env.KafkaBrokerUrl, cfg.Env.KafkaTopicManagerToWorkers, cfg.Env.KafkaTopicWorkersToManager)
	if err != nil {
		return nil, tracerr.Wrap(err)
	}
	kafkaSender := myKafka.New(kafkaWriter)

	// Подключение к MongoDB
	mongoClient, err := connectMongo(cfg.Env.MongoUri)
	if err != nil {
		return nil, tracerr.Wrap(err)
	}
	mongoDb, err := initMongo(mongoClient, cfg.Env.IsDebug)
	if err != nil {
		return nil, tracerr.Wrap(err)
	}

	repo := storage.New(mongoDb)
	scheduler := sch.CreateScheduler(repo, kafkaSender)
	monitor := monitor.CreateMonitor(scheduler)

	crack := usecases.Crack{
		ManagerRepo: repo,
		Scheduler:   scheduler,
		Monitor:     monitor,
	}

	server := httpserver.NewServer(&crack)
	consumer := kafka_consumer.New(kafkaReader, &crack)
	return &App{
		config:        cfg,
		server:        &server,
		scheduler:     scheduler,
		kafkaConsumer: consumer,
		kafkaSender:   kafkaSender,
	}, nil
}

func (a *App) StartApp() {
	a.scheduler.Run()
	a.kafkaConsumer.Run()
	log.Println("HttpServer завершил работу: ", a.server.ListenAndServe())
}

func connectKafka(brokerURL, managerToWorkersTopic, workersToManagerTopic string) (*kafka.Writer, *kafka.Reader, error) {
	// Настройка Kafka Writer
	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers:      []string{brokerURL},   // Адрес Kafka-брокера
		Topic:        managerToWorkersTopic, // Тема для отправки сообщений от менеджера к воркерам
		Balancer:     &kafka.LeastBytes{},   // Балансировщик для выбора партиции
		RequiredAcks: int(kafka.RequireAll),
		BatchSize:    1,
	})

	// Настройка Kafka Reader
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:        []string{brokerURL},   // Адрес Kafka-брокера
		Topic:          workersToManagerTopic, // Тема для получения сообщений от воркеров
		GroupID:        "manager-group",       // ID группы потребителей
		CommitInterval: 0,
		StartOffset:    kafka.LastOffset,
	})

	log.Println("Connected to Kafka successfully")
	return writer, reader, nil
}

func connectMongo(mongoURI string) (*mongo.Client, error) {
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(mongoURI))
	if err != nil {
		return nil, tracerr.New(fmt.Sprintf("failed to create MongoDB client: %v", err))
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, tracerr.New(fmt.Sprintf("failed to connect to MongoDB: %v", err))
	}
	log.Println("Connected to MongoDB successfully")
	return client, nil
}

func initMongo(client *mongo.Client, isDebug bool) (*mongo.Database, error) {
	if isDebug {
		if err := client.Database("crack").Drop(context.Background()); err != nil {
			return nil, tracerr.Wrap(err)
		}
	}

	db := client.Database("crack")
	collections, err := db.ListCollectionNames(context.Background(), bson.D{})
	if err != nil {
		return nil, tracerr.Wrap(err)
	}

	if len(collections) > 0 && !isDebug {
		return db, nil
	}

	//init
	err = db.CreateCollection(context.Background(), "seq")
	if err != nil {
		return nil, tracerr.Wrap(err)
	}

	err = db.CreateCollection(context.Background(), "orders")
	if err != nil {
		return nil, tracerr.Wrap(err)
	}

	err = db.CreateCollection(context.Background(), "tasks")
	if err != nil {
		return nil, tracerr.Wrap(err)
	}

	err = db.CreateCollection(context.Background(), "workers")
	if err != nil {
		return nil, tracerr.Wrap(err)
	}

	// Создание индекса для поля order_id
	taskCollection := db.Collection("tasks")
	_, err = taskCollection.Indexes().CreateOne(context.Background(), mongo.IndexModel{
		Keys:    bson.D{{Key: "order_id", Value: 1}}, // Индекс по полю order_id
		Options: options.Index().SetUnique(false),    // Уникальность можно включить, если нужно
	})
	if err != nil {
		return nil, tracerr.Wrap(err)
	}

	_, err = taskCollection.Indexes().CreateOne(context.Background(), mongo.IndexModel{
		Keys: bson.D{
			{Key: "order_id", Value: 1},
			{Key: "block_number", Value: 1},
		},
		Options: options.Index().SetUnique(true),
	})
	if err != nil {
		return nil, tracerr.Wrap(err)
	}

	return db, nil
}
