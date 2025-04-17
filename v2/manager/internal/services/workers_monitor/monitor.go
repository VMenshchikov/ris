package monitor

import (
	"context"
	"hash_manager/internal/services/scheduler"
	"log"
	"sync"
	"time"
)

const (
	KEEP_ALIVE_INTERVAL = 10 * time.Second
)

type Monitor struct {
	scheduler *scheduler.Scheduler
	workers   map[uint64]*workerInfo
	mu        sync.Mutex
}

type workerInfo struct {
	id         uint64
	maxTasks   uint
	lastActive time.Time
	cancel     context.CancelFunc
}

func CreateMonitor(scheduler *scheduler.Scheduler) *Monitor {
	return &Monitor{
		scheduler: scheduler,
		workers:   make(map[uint64]*workerInfo),
	}
}

func (m *Monitor) UpdateWorker(id uint64, maxTasks uint) {
	m.mu.Lock()
	defer m.mu.Unlock()

	worker, exists := m.workers[id]
	if exists {
		log.Println("Обновляем воркера", id)
		worker.cancel()
	} else {
		log.Println("Создаем воркера", id)
		worker = &workerInfo{
			id:       id,
			maxTasks: maxTasks,
		}
		m.workers[id] = worker
		m.scheduler.UpdateMaxTasks(int64(worker.maxTasks))
	}

	ctx, cancel := context.WithTimeout(context.Background(), KEEP_ALIVE_INTERVAL)
	worker.cancel = cancel
	worker.lastActive = time.Now()

	go func(w *workerInfo) {
		<-ctx.Done()
		if ctx.Err() == context.DeadlineExceeded {
			m.mu.Lock()
			defer m.mu.Unlock()

			delete(m.workers, w.id)
			m.scheduler.UpdateMaxTasks(-int64(w.maxTasks))
			log.Println("Удалили воркера", id)
		}
	}(worker)
}
