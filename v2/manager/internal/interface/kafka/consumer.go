package kafka_consumer

import (
	"context"
	"encoding/json"
	"hash_manager/internal/interface/kafka/dto"
	"hash_manager/internal/usecases"
	"log"

	"github.com/segmentio/kafka-go"
	"github.com/ztrue/tracerr"
)

// KafkaHandler - структура для обработки сообщений из Kafka
type KafkaConsumer struct {
	reader *kafka.Reader
	crack  *usecases.Crack
}

func New(reader *kafka.Reader, crack *usecases.Crack) *KafkaConsumer {
	return &KafkaConsumer{
		reader: reader,
		crack:  crack,
	}
}

func (h *KafkaConsumer) Run() {
	go h.processMessages()
}

// ProcessMessages - метод для обработки сообщений из Kafka
func (h *KafkaConsumer) processMessages() {
	for {
		msg, err := h.reader.FetchMessage(context.Background())
		if err != nil {
			log.Printf("Ошибка при чтении сообщения: %v", err)
			//return
			continue
		}

		switch string(msg.Key) {
		case "alive":
			if err := h.handleAlive(msg); err != nil {
				log.Println(err)
				continue
			}
		case "result":
			if err := h.handleResult(msg); err != nil {
				log.Println(err)
				continue
			}

		default:
			log.Printf("Unknown message type: %s", msg.Key)
		}

		err = h.reader.CommitMessages(context.Background(), msg)
		if err != nil {
			log.Printf("Ошибка при подтверждении сообщения: %v", err)
		}
	}
}

func (h *KafkaConsumer) Close() {
	if err := h.reader.Close(); err != nil {
		log.Printf("Ошибка при закрытии Reader: %v", err)
	}
}

func (h *KafkaConsumer) handleAlive(msg kafka.Message) error {
	var data dto.AliveMessage
	if err := json.Unmarshal(msg.Value, &data); err != nil {
		return tracerr.Wrap(err)
	}

	h.crack.UpdateWorker(data.WorkerId, data.MaxTasks)
	return nil
}

func (h *KafkaConsumer) handleResult(msg kafka.Message) error {
	var data dto.ResultMessage
	if err := json.Unmarshal(msg.Value, &data); err != nil {
		return tracerr.Wrap(err)
	}
	if err := h.crack.EndTask(data.OrderId, data.TaskNumber, data.Results...); err != nil {
		return tracerr.Wrap(err)
	}
	return nil

}
