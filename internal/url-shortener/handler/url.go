package handler

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"net/http"
)

func (h *Handler) SaveURL(w http.ResponseWriter, r *http.Request) {
	var req SaveURLRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	if err := req.validate(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.urlService.URLSave(req.URL, req.Alias); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(SaveURLResponse{
		Status: "OK",
	})
}

func (h *Handler) GetURL(w http.ResponseWriter, r *http.Request) {
	alias := chi.URLParam(r, "alias")

	if alias == "" {
		http.Error(w, "alias is required", http.StatusBadRequest)
		return
	}

	url, err := h.urlService.URLGet(alias)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(GetURLResponse{
		URL: url,
	})
}
