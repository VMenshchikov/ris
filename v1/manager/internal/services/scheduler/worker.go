package scheduler

import (
	"hash_manager/internal/client"
	"log"
	"sync"
	"sync/atomic"

	"github.com/ztrue/tracerr"
)

type Worker struct {
	id           uint
	scheduler    *Scheduler
	urn          string
	maxTasks     uint
	currentTasks *atomic.Int32

	lock sync.Mutex
}

func CreateWorker(scheduler *Scheduler, urn string, maxTasks, id uint) *Worker {
	return &Worker{
		id:           id,
		scheduler:    scheduler,
		urn:          urn,
		maxTasks:     maxTasks,
		currentTasks: &atomic.Int32{},
	}
}

func (d *Worker) GetScheduler() *Scheduler {
	return d.scheduler
}

func (d *Worker) CompleteTask() {
	d.currentTasks.Add(-1)
	d.sendTask()
}

func (d *Worker) Run() {
	log.Println("Запуск воркера:", d.id)
	d.sendTask()
}

func (d *Worker) sendTask() {
	d.lock.Lock()

	if uint(d.currentTasks.Load()) == d.maxTasks {
		d.lock.Unlock()
		return
	}
	task := <-d.scheduler.GetTaskChannel()
	err := client.SendTask(d.urn, client.TaskDto{
		OrderId:     task.OrderId,
		TargetHash:  task.TargetHash,
		BlockSize:   task.BlockSize,
		BlockNumber: task.BlockNumber,
		MaxLen:      task.MaxLen,
	})

	if err != nil {
		tracerr.Print(err)
		d.scheduler.GetTaskChannel() <- task
		d.lock.Unlock()
		return
	}

	d.currentTasks.Add(1)
	if uint(d.currentTasks.Load()) < d.maxTasks {
		d.lock.Unlock()
		d.sendTask()
	} else {
		d.lock.Unlock()
	}

}
