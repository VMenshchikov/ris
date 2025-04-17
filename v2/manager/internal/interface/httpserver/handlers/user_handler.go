package handlers

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"hash_manager/internal/interface/httpserver/dto"
	"hash_manager/internal/usecases"
	"net/http"
	"regexp"
	"strconv"
)

type userHandler struct {
	crack *usecases.Crack
}

func CreateUserHandler(crack *usecases.Crack) *userHandler {
	return &userHandler{
		crack: crack,
	}
}

func (h *userHandler) CrackHashHandle(w http.ResponseWriter, r *http.Request) {
	var data dto.CrackHashRequest
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil || !data.SetDefaults() || !isMd5Hash(data.Hash) {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}

	ArrayHash, err := hex.DecodeString(data.Hash)
	if err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}
	var fixedArrayHash [16]byte
	copy(fixedArrayHash[:], ArrayHash)

	id, err := h.crack.CreateOrder(
		fixedArrayHash,
		data.MaxLength,
		data.Timeout,
		data.BlockSize,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	jsonResponse, err := json.Marshal(dto.CrackHashResponse{
		Id: id,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusAccepted)
	w.Write(jsonResponse)
}

func (h *userHandler) GetResultHandle(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(r.URL.Query()["id"][0], 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	status, results, percentage, err := h.crack.GetResult(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	jsonResponse, err := json.Marshal(dto.GetResultResponse{
		Status:   status,
		Results:  results,
		Progress: fmt.Sprintf("%.02f", percentage),
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}

func isMd5Hash(value string) bool {
	re := regexp.MustCompile("^[a-fA-F0-9]{32}$")
	return re.MatchString(value)
}
