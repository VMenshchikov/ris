package kafka

import (
	"context"
	"encoding/json"
	"fmt"

	gokafka "github.com/segmentio/kafka-go"
	"github.com/ztrue/tracerr"
)

type MessageTask struct {
	OrderId     uint64
	TargetHash  [16]byte
	BlockSize   uint
	BlockNumber uint
	MaxLen      uint
}

type Sender struct {
	writer *gokafka.Writer
}

func New(writer *gokafka.Writer) *Sender {
	return &Sender{
		writer: writer,
	}
}

func (s *Sender) Send(data MessageTask) error {

	json, err := json.Marshal(data)
	if err != nil {
		return tracerr.New(fmt.Sprintf("failed to write message: %v", err))
	}

	message := gokafka.Message{
		Key:   nil,
		Value: json,
	}

	err = s.writer.WriteMessages(context.Background(), message)
	if err != nil {
		return tracerr.New(fmt.Sprintf("failed to write message: %v", err))
	}
	return nil
}
