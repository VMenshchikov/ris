package app

import (
	"hash_worker/internal/interface/httpserver"
	"hash_worker/internal/usecases"
	"log"
	"runtime"
)

type App struct {
	config Config
	server *httpserver.Server
}

func NewApp(cfg Config) App {
	crack := usecases.Crack{
		WorkerID:       cfg.WorkerId,
		ManagerAddress: cfg.ManagerAddress,
	}

	server := httpserver.NewServer(&crack)
	return App{
		config: cfg,
		server: &server,
	}
}

func (a *App) Initialize() {
	if a.config.MaxProcs < uint(runtime.NumCPU()) {
		runtime.GOMAXPROCS(int(a.config.MaxProcs))
	}
	log.Printf("Procs: %d/%d", a.config.MaxProcs, runtime.NumCPU())
}

func (a *App) RunApp() {
	log.Println("Запуск сервера")
	a.server.ListenAndServe()
	log.Println("Сервер запущен")
}
