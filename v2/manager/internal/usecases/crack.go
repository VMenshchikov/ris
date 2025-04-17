package usecases

import (
	"hash_manager/internal/domain/model"
	"hash_manager/internal/domain/repo"
	"hash_manager/internal/services/scheduler"
	monitor "hash_manager/internal/services/workers_monitor"
	"log"
	"math"
	"time"

	"github.com/ztrue/tracerr"
)

type Crack struct {
	ManagerRepo repo.ManagerRepository
	Scheduler   *scheduler.Scheduler
	Monitor     *monitor.Monitor
}

func (c *Crack) CreateOrder(targetHash [16]byte, maxLen uint, timeout uint, blockSize uint) (uint64, error) {

	timeoutTime := time.Now().Add(time.Duration(timeout) * time.Second)

	id, err := c.ManagerRepo.AddOrder(targetHash, maxLen, timeoutTime, blockSize)
	if err != nil {
		return 0, tracerr.Wrap(err)
	}

	return id, nil
}

func (c *Crack) GetResult(id uint64) (string, []string, float64, error) {
	order, err := c.ManagerRepo.FindOrder(id)
	if err != nil {
		return "", nil, 0, tracerr.Wrap(err)
	}

	completedTasksCount, err := c.ManagerRepo.CountCompletedTasksByOrderId(order.Id)
	if err != nil {
		return "", nil, 0, tracerr.Wrap(err)
	}

	allTasksCount := uint(math.Ceil(math.Pow(float64(36), float64(order.MaxLen)) / float64(order.BlockSize)))

	return string(order.Status), order.Results, float64(completedTasksCount) / float64(allTasksCount) * 100, nil
}

func (c *Crack) FinishOrder(id uint64, status model.OrderStatusType) error {
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

func (c *Crack) EndTask(orderId uint64, blockNumber uint, results ...string) error {
	log.Println("Получен ответ от воркера. Заказ", orderId, ", блок:", blockNumber, ", результаты: ", results)

	if len(results) != 0 {
		order, err := c.ManagerRepo.FindOrder(orderId)
		if err != nil {
			return tracerr.Wrap(err)
		}
		order.Results = append(order.Results, results...)
		c.ManagerRepo.UpdateOrder(order)
	}

	task, err := c.ManagerRepo.FindTask(orderId, uint64(blockNumber))
	if err != nil {
		return tracerr.Wrap(err)
	}

	task.Status = model.COMPLETED
	c.ManagerRepo.UpdateTask(*task)

	completedTasksCount, err := c.ManagerRepo.CountCompletedTasksByOrderId(orderId)
	if err != nil {
		return tracerr.Wrap(err)
	}
	order, err := c.ManagerRepo.FindOrder(task.OrderId)
	if err != nil {
		return tracerr.Wrap(err)
	}

	if completedTasksCount == int64(math.Ceil(math.Pow(float64(36), float64(order.MaxLen))/float64(order.BlockSize))) {
		order.Status = model.OrderStatusType(model.COMPLETED)
		if err := c.ManagerRepo.UpdateOrder(order); err != nil {
			return tracerr.Wrap(err)
		}
	}

	c.Scheduler.FinishTask()
	return nil
}

func (c *Crack) UpdateWorker(workerId uint, maxTasks uint) {
	c.Monitor.UpdateWorker(uint64(workerId), maxTasks)
}
