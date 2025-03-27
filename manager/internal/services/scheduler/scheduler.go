package scheduler

import (
	"hash_manager/internal/domain/model"
	"hash_manager/internal/domain/repo"
	"log"
	"math"
	"time"
)

type Scheduler struct {
	orders       chan scheduledOrder
	tasksChannel chan Task

	repo repo.ManagerRepository
}

type scheduledOrder struct {
	OrderId    uint64
	Weight     uint
	TargetHash [16]byte
	BlockSize  uint
	MaxLen     uint

	CurrentBlock uint
	MaxBlocks    uint

	closeChannel chan interface{}
}

type Task struct {
	OrderId     uint64
	TargetHash  [16]byte
	BlockSize   uint
	BlockNumber uint
	MaxLen      uint
}

func CreateScheduler(repo repo.ManagerRepository) Scheduler {
	return Scheduler{
		orders:       make(chan scheduledOrder, 5),
		tasksChannel: make(chan Task, 2),
		repo:         repo,
	}
}

func (s *Scheduler) Run() {
	go func() {
		for {
			order := <-s.orders

			select {
			case <-order.closeChannel:
				log.Println("Шедулер выкинул просроченный заказ:", order.OrderId)

				oModel, _ := s.repo.FindOrder(order.OrderId)
				oModel.Lock.Lock()
				if oModel.Status == model.IN_PROGRESS {
					oModel.Status = model.ERROR
				}
				oModel.Lock.Unlock()
				continue //канал закрыт, выбрасываем из очереди
			default:
			}

			for i := 0; i < int(order.Weight) && order.CurrentBlock != order.MaxBlocks; i++ {
				log.Printf("Шедулер создает задачу %d/%d для заказа %d", order.CurrentBlock+1, order.MaxBlocks, order.OrderId)

				s.tasksChannel <- Task{
					OrderId:     uint64(order.OrderId),
					TargetHash:  order.TargetHash,
					BlockSize:   order.BlockSize,
					BlockNumber: order.CurrentBlock,
					MaxLen:      order.MaxLen,
				}
				order.CurrentBlock++
			}
			if order.CurrentBlock == order.MaxBlocks {
				close(order.closeChannel)
				mOrder, _ := s.repo.FindOrder(order.OrderId)

				mOrder.Lock.Lock()
				mOrder.Status = model.ALL_TASKS
				mOrder.Lock.Unlock()

			} else {
				s.orders <- order
			}
		}
	}()
}

func (s *Scheduler) AddOrder(order scheduledOrder) {
	s.orders <- order
}

func (s *Scheduler) GetTaskChannel() chan Task {
	return s.tasksChannel
}

func CreateScheduledOrder(orderId uint64, targetHash [16]byte, blockSize uint, maxLen uint, weight uint, timeout time.Time) *scheduledOrder {
	var order = &scheduledOrder{
		OrderId:      orderId,
		Weight:       weight,
		TargetHash:   targetHash,
		BlockSize:    blockSize,
		MaxLen:       maxLen,
		CurrentBlock: 0,
		MaxBlocks:    uint(math.Ceil(float64((math.Pow(36, float64(maxLen)))) / float64(blockSize))),
		closeChannel: make(chan interface{}),
	}
	go order.timeoutObserver(timeout)
	return order
}

func (o *scheduledOrder) timeoutObserver(timeout time.Time) {
	timer := time.NewTimer(time.Until(timeout))
	defer timer.Stop()

	select {
	case <-timer.C:
		close(o.closeChannel)
	case <-o.closeChannel:
	}
}
