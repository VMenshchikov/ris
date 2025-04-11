package handlers

import (
	"encoding/json"
	"hash_manager/internal/interface/httpserver/dto"
	"hash_manager/internal/usecases"
	"net/http"
)

type workerHandler struct {
	crack *usecases.Crack
}

func CreateWorkerHandler(crack *usecases.Crack) *workerHandler {
	return &workerHandler{
		crack: crack,
	}
}

func (h *workerHandler) WorkerResponseHandle(w http.ResponseWriter, r *http.Request) {
	var data dto.WorkerResult
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	w.WriteHeader(http.StatusAccepted)
	h.crack.EndTask(data.WorkerId, data.OrderId, data.Results...)
}
