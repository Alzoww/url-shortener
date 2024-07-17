package handler

import (
	"github.com/Alzoww/url-shortener/internal/storage"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Handler struct {
	storage storage.Interface
}

func New(storage storage.Interface) *Handler {
	return &Handler{
		storage: storage,
	}
}

func (h *Handler) InitRoutes() *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Post("/v1/url", h.SaveURL)
	r.Get("/v1/url/{alias}", h.GetURL)
	r.Get("/v1/{alias}", h.Redirect)

	return r
}
