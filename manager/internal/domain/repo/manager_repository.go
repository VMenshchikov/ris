package repo

import (
	"hash_manager/internal/domain/model"
	"time"
)

type ManagerRepository interface {
	AddOrder(targetHash [16]byte, maxLen uint, timeout time.Time, weight uint, blockSize uint) (uint64, error)
	UpdateOrder(*model.OrderInfo) error
	FindOrder(uint64) (*model.OrderInfo, error)
	/*
		FindResults(orderId int) (error, []string)
	*/
}
