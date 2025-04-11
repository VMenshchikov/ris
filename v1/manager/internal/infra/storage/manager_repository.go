package storage

import (
	"hash_manager/internal/domain/model"
	"hash_manager/internal/domain/repo"
	"math"
	"sync"
	"time"

	"github.com/ztrue/tracerr"
)

var (
	_ repo.ManagerRepository = (*ManagerRepository)(nil)
)

type ManagerRepository struct {
	orders map[uint64]*model.OrderInfo
	count  uint64
}

func New() *ManagerRepository {
	return &ManagerRepository{
		orders: make(map[uint64]*model.OrderInfo, 50),
		count:  0,
	}
}

func (m *ManagerRepository) AddOrder(targetHash [16]byte, maxLen uint, timeout time.Time, weight uint, blockSize uint) (uint64, error) {
	m.count++
	m.orders[m.count] = &model.OrderInfo{
		Id:               m.count,
		Status:           model.IN_PROGRESS,
		TargetHash:       targetHash,
		MaxLen:           maxLen,
		Timeout:          timeout,
		Weight:           weight,
		BlockSize:        int64(blockSize),
		Results:          []string{},
		TasksUncompleted: uint(math.Pow(float64(36), float64(maxLen)) / float64(blockSize)),
		Lock:             sync.Mutex{},
	}
	return m.count, nil
}

func (m *ManagerRepository) UpdateOrder(order *model.OrderInfo) error {
	_, contains := m.orders[order.Id]
	if !contains {
		return tracerr.New("no present")
	}

	m.orders[order.Id] = order
	return nil
}

func (m *ManagerRepository) FindOrder(id uint64) (*model.OrderInfo, error) {
	value, contains := m.orders[id]
	if !contains {
		return &model.OrderInfo{}, tracerr.New("no present")
	}

	return value, nil
}
