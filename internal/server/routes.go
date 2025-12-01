package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (s *Server) routes() {
	r := chi.NewRouter()

	// Healthcheck endpoint
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	})

	// Buraya sonra:
	// - /requests
	// - /requests/{id}
	// - /requests/{id}/approve
	// - /requests/{id}/reject
	// gibi route'ları ekleyeceğiz.

	s.router = r
}
