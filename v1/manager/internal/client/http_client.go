package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ztrue/tracerr"
)

const (
	workerApi = "/internal/api/worker/hash/crack/task"
)

type TaskDto struct {
	OrderId     uint64   `json:"orderId"`
	TargetHash  [16]byte `json:"targetHash"`
	BlockSize   uint     `json:"blockSize"`
	BlockNumber uint     `json:"blockNumber"`
	MaxLen      uint     `json:"maxLen"`
}

func SendTask(urn string, task TaskDto) error {
	client := &http.Client{}
	value, err := json.Marshal(task)
	if err != nil {
		return tracerr.New(fmt.Sprintln("Не удалось создать json: ", task))
	}

	req, err := http.NewRequest("POST", "http://"+urn+workerApi, bytes.NewBuffer(value))
	if err != nil {
		return tracerr.New(fmt.Sprintln("Ошибка создания запроса:", err))
	}

	req.Header.Set("Content-Type", "application/json")

	_, err = client.Do(req)
	if err != nil {
		return tracerr.New(fmt.Sprintln("Ошибка выполнения запроса:", err))
	}

	return nil
}
