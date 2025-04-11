package handlers

import (
	"encoding/json"
	"hash_worker/internal/domain"
	"hash_worker/internal/interface/httpserver/dto"
	"hash_worker/internal/usecases"
	"log"
	"net/http"
)

type workerHandler struct {
	crack *usecases.Crack
}

func CreateWorkerHandle(crack *usecases.Crack) *workerHandler {
	return &workerHandler{
		crack: crack,
	}
}

func (h *workerHandler) NewTaskHandle(w http.ResponseWriter, r *http.Request) {
	log.Println("Получена задача")
	var data dto.Task
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}

	h.crack.CrackHash(domain.Task{
		OrderId:     uint(data.OrderId),
		TargetHash:  data.TargetHash,
		MaxLen:      data.MaxLen,
		BlockNumber: data.BlockNumber,
		BlockSize:   data.BlockSize,
	})
}
