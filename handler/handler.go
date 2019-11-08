package handler

import (
	"blendcube/conf"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func New() http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(conf.Shared.Server.Timeout * time.Second))

	r.Get("/status", statusCheck)
	r.Get("/cube", generateCube)
	r.Get("/cube.gltf", generateCube)
	r.Get("/cube.glb", generateCube)

	return r
}

func statusCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func generateCube(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
