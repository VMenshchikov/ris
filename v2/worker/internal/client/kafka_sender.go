package client

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/segmentio/kafka-go"
	"github.com/ztrue/tracerr"
)

type TaskDto struct {
	WorkerId   uint     `json:"workerId"`
	OrderId    uint64   `json:"orderId"`
	TaskNumber uint     `json:"number"`
	Results    []string `json:"results"`
}

type AliveDto struct {
	WorkerId uint `json:"workerId"`
	MaxTasks uint `json:"maxTasks"`
}

type Sender struct {
	writer *kafka.Writer
}

func New(writer *kafka.Writer) *Sender {
	return &Sender{
		writer: writer,
	}
}

func (s *Sender) SendResult(data TaskDto) error {
	json, err := json.Marshal(data)
	if err != nil {
		return tracerr.New(fmt.Sprintf("failed to write message: %v", err))
	}

	message := kafka.Message{
		Key:   []byte("result"),
		Value: json,
	}
	log.Println("start send", time.Now())
	err = s.writer.WriteMessages(context.Background(), message)
	if err != nil {
		tracerr.New(fmt.Sprintf("failed to write message: %v", err))
	}
	log.Println("end send", time.Now())
	return nil
}

func (s *Sender) SendAlive(data AliveDto) error {
	json, err := json.Marshal(data)
	if err != nil {
		return tracerr.New(fmt.Sprintf("failed to write message: %v", err))
	}

	message := kafka.Message{
		Key:   []byte("alive"),
		Value: json,
	}

	err = s.writer.WriteMessages(context.Background(), message)
	if err != nil {
		tracerr.New(fmt.Sprintf("failed to write message: %v", err))
	}
	return nil
}
