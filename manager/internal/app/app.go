package app

import (
	"hash_manager/internal/infra/storage"
	"hash_manager/internal/interface/httpserver"
	sch "hash_manager/internal/services/scheduler"
	"hash_manager/internal/usecases"
	"log"
)

type App struct {
	config    Config
	server    *httpserver.Server
	scheduler sch.Scheduler
	workers   map[uint]*sch.Worker
}

func NewApp(cfg Config) App {

	repo := storage.New()
	scheduler := sch.CreateScheduler(repo)

	workers := make(map[uint]*sch.Worker)
	for i := 0; i < len(cfg.WorkersId); i++ {
		workers[uint(i+1)] = sch.CreateWorker(&scheduler, cfg.WorkerAddresses[i], cfg.MaxTasks[i], cfg.WorkersId[i])
	}

	crack := usecases.Crack{
		ManagerRepo: repo,
		Scheduler:   &scheduler,
		Workers:     workers,
	}

	server := httpserver.NewServer(&crack)

	return App{
		config:    cfg,
		server:    &server,
		scheduler: scheduler,
		workers:   workers,
	}
}

func (a *App) StartApp() {
	a.scheduler.Run()
	for _, v := range a.workers {
		go v.Run()
	}
	log.Println("HttpServer завершил работу: ", a.server.ListenAndServe())
}
