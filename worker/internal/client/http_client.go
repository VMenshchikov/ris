package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ztrue/tracerr"
)

const (
	managerApi = "/internal/api/hash/"
)

type TaskDto struct {
	WorkerId uint     `json:"workerId"`
	OrderId  uint64   `json:"orderId"`
	Results  []string `json:"results"`
}

func SendResult(urn string, task TaskDto) error {
	client := &http.Client{}
	value, err := json.Marshal(task)
	if err != nil {
		return tracerr.New(fmt.Sprintln("Не удалось создать json: ", task))
	}

	req, err := http.NewRequest("PATCH", "http://"+urn+managerApi, bytes.NewBuffer(value))
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
