package httpserver

import (
	"hash_manager/internal/interface/httpserver/handlers"
	"hash_manager/internal/usecases"
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
)

type Server struct {
	server http.Server
}

func initMiddlewares(router *chi.Mux) {
	//router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(render.SetContentType(render.ContentTypeJSON))
}

func initRoutes(router *chi.Mux, crack *usecases.Crack) {

	userHandler := handlers.CreateUserHandler(crack)
	//workerHandler := handlers.CreateWorkerHandler(crack)

	router.Get("/*", func(w http.ResponseWriter, r *http.Request) {
		log.Println("!!!!")
	})

	router.Post("/api/hash/crack", userHandler.CrackHashHandle)
	router.Get("/api/hash/status", userHandler.GetResultHandle)

	//router.Patch("/internal/api/hash/", workerHandler.WorkerResponseHandle)

}

func NewServer(crack *usecases.Crack) Server {
	r := chi.NewRouter()

	initMiddlewares(r)
	initRoutes(r, crack)

	return Server{
		server: http.Server{
			Addr:    ":8081",
			Handler: r,
		},
	}
}

func (s *Server) ListenAndServe() error {
	log.Println("Начали слушать")
	return s.server.ListenAndServe()
}
