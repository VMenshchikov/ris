package kafka

import (
	"context"
	"encoding/json"
	"hash_worker/internal/domain"
	"hash_worker/internal/usecases"
	"log"
	"time"

	"github.com/segmentio/kafka-go"
	"github.com/ztrue/tracerr"
)

type Consumer struct {
	reader *kafka.Reader
	crack  *usecases.Crack
}

func New(reader *kafka.Reader, crack *usecases.Crack) *Consumer {
	return &Consumer{
		reader: reader,
		crack:  crack,
	}
}

func (c *Consumer) Run() {
	c.processMessages()
}

func (c *Consumer) processMessages() {
	for {
		log.Println("start task", time.Now())
		msg, err := c.reader.FetchMessage(context.Background())
		if err != nil {
			log.Printf("Ошибка при чтении сообщения: %v", err)
			continue
		}
		if err := c.handleTask(msg); err != nil {
			log.Println(err)
			continue
		}

		err = c.reader.CommitMessages(context.Background(), msg)
		if err != nil {
			log.Printf("Ошибка при подтверждении сообщения: %v", err)
		}
		log.Println("end task", time.Now())
	}
}

func (c *Consumer) handleTask(msg kafka.Message) error {
	var task domain.Task
	if err := json.Unmarshal(msg.Value, &task); err != nil {
		return tracerr.Wrap(err)
	}

	return c.crack.CrackHash(task)
}
