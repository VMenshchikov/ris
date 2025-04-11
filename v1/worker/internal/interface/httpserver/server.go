package httpserver

import (
	"hash_worker/internal/interface/httpserver/handlers"
	"hash_worker/internal/usecases"
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
	userHandler := handlers.CreateWorkerHandle(crack)

	router.Post("/internal/api/worker/hash/crack/task", userHandler.NewTaskHandle)
}

func NewServer(crack *usecases.Crack) Server {
	r := chi.NewRouter()

	initMiddlewares(r)
	initRoutes(r, crack)

	return Server{
		server: http.Server{
			Addr:    ":8082",
			Handler: r,
		},
	}
}

func (s *Server) ListenAndServe() error {
	log.Println("Начали слушать")
	return s.server.ListenAndServe()
}
