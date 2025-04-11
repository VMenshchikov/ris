package usecases

import (
	"hash_manager/internal/domain/model"
	"hash_manager/internal/domain/repo"
	"hash_manager/internal/services/scheduler"
	"log"
	"math"
	"time"

	"github.com/ztrue/tracerr"
)

type Crack struct {
	ManagerRepo repo.ManagerRepository
	Scheduler   *scheduler.Scheduler
	Workers     map[uint]*scheduler.Worker
}

func (c *Crack) CreateOrder(targetHash [16]byte, maxLen uint, timeout uint, weight uint, blockSize uint) (uint64, error) {

	timeoutTime := time.Now().Add(time.Duration(timeout) * time.Second)

	id, err := c.ManagerRepo.AddOrder(targetHash, maxLen, timeoutTime, weight, blockSize)
	if err != nil {
		return 0, tracerr.Wrap(err)
	}

	c.Scheduler.AddOrder(*scheduler.CreateScheduledOrder(id, targetHash, blockSize, maxLen, weight, timeoutTime))

	return id, nil
}

func (c *Crack) GetResult(id uint64) (string, []string, float64, error) {
	order, err := c.ManagerRepo.FindOrder(id)
	if err != nil {
		return "", nil, 0, tracerr.Wrap(err)
	}

	allTasks := uint(math.Pow(float64(36), float64(order.MaxLen)) / float64(order.BlockSize))

	return string(order.Status), order.Results, float64(allTasks-order.TasksUncompleted) / float64(allTasks) * 100, nil
}

func (c *Crack) FinishOrder(id uint64, status model.StatusType) error {
	order, err := c.ManagerRepo.FindOrder(id)
	if err != nil {
		return tracerr.Wrap(err)
	}

	order.Status = status
	err = c.ManagerRepo.UpdateOrder(order)
	if err != nil {
		return tracerr.Wrap(err)
	}

	return nil
}

func (c *Crack) EndTask(workerId uint, orderId uint64, results ...string) error {
	go c.Workers[workerId].CompleteTask()

	log.Println("Получен ответ от воркера", workerId, ". Pаказ", orderId, ", результаты: ", results)

	order, err := c.ManagerRepo.FindOrder(orderId)
	if err != nil {
		return tracerr.Wrap(err)
	}

	order.Lock.Lock()
	defer order.Lock.Unlock()

	order.TasksUncompleted--
	order.Results = append(order.Results, results...)
	if order.TasksUncompleted == 0 {
		order.Status = model.READY
	}
	return nil
}
