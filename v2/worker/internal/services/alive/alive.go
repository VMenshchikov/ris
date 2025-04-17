package alive

import (
	"hash_worker/internal/client"
	"time"
)

type AliveService struct {
	Sender   *client.Sender
	WorkerId uint
	MaxTasks uint
}

const ALIVE_MSG_INTERVAL = time.Second * 3

func (a *AliveService) TranslateAlive() {
	for {
		a.Sender.SendAlive(client.AliveDto{
			WorkerId: a.WorkerId,
			MaxTasks: a.MaxTasks,
		})
		time.Sleep(ALIVE_MSG_INTERVAL)
	}
}
