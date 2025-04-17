package repo

import (
	"hash_manager/internal/domain/model"
	"time"
)

type ManagerRepository interface {
	AddOrder(targetHash [16]byte, maxLen uint, timeout time.Time, blockSize uint) (uint64, error)
	UpdateOrder(*model.OrderInfo) error
	FindOrder(uint64) (*model.OrderInfo, error)
	FindOrderForExecution() (*model.OrderInfo, error)
	CloseTimeoutOrders() error

	AddTask(task model.TaskInfo) error
	UpdateTask(task model.TaskInfo) error
	FindTasksByOrderId(orderId uint64) ([]model.TaskInfo, error)
	FindTask(orderId, blockNumber uint64) (*model.TaskInfo, error)
	CountCompletedTasksByOrderId(orderId uint64) (int64, error)
	FindOutdatedSendedTasks(before time.Time) ([]model.TaskInfo, error)

	//FindProgress()

	AddWorker(worker model.Worker) error
	UpdateWorker(worker model.Worker) error
	FindWorker(id int64) (model.Worker, error)
}
