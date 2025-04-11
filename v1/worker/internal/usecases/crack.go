package usecases

import (
	"hash_worker/internal/client"
	"hash_worker/internal/domain"
	"hash_worker/internal/services"
	"log"
)

type Crack struct {
	WorkerID       uint
	ManagerAddress string
}

func (c *Crack) CrackHash(task domain.Task) {
	go func() {
		log.Printf("Начали взлом c %d по %d", task.BlockSize*task.BlockNumber, task.BlockSize*task.BlockNumber+task.BlockSize)
		res := services.CheckRange(task.TargetHash, task.MaxLen, task.BlockSize, uint64(task.BlockSize*task.BlockNumber))
		log.Println("Результаты:", res)
		client.SendResult(c.ManagerAddress, client.TaskDto{
			WorkerId: c.WorkerID,
			OrderId:  uint64(task.OrderId),
			Results:  res,
		})
	}()
}
