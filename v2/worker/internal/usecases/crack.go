package usecases

import (
	"hash_worker/internal/client"
	"hash_worker/internal/domain"
	"hash_worker/internal/services/brutforce"
	"log"
)

type Crack struct {
	WorkerID uint
	Sender   *client.Sender
}

func (c *Crack) CrackHash(task domain.Task) error {
	log.Printf("Начали взлом c %d по %d", task.BlockSize*task.BlockNumber, task.BlockSize*task.BlockNumber+task.BlockSize)
	res := brutforce.CheckRange(task.TargetHash, task.MaxLen, task.BlockSize, uint64(task.BlockSize*task.BlockNumber))
	log.Println("Результаты:", res)
	return c.Sender.SendResult(client.TaskDto{
		WorkerId:   c.WorkerID,
		OrderId:    uint64(task.OrderId),
		TaskNumber: task.BlockNumber,
		Results:    res,
	})
}
