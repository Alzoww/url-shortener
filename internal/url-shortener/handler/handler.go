package handler

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type URLServiceI interface {
	URLSave(urlToSave, alias string) error
	URLGet(alias string) (string, error)
}

type Handler struct {
	urlService URLServiceI
}

func New(urlService URLServiceI) *Handler {
	return &Handler{
		urlService: urlService,
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
