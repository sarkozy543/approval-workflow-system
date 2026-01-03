package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

func (s *Server) routes() {
	r := chi.NewRouter()

	// ✅ TÜM MIDDLEWARE'LER BURADA, ROUTE'LARDAN ÖNCE
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173"}, // frontend dev server
		AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Content-Type", "X-User", "Origin"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// ✅ HEALTH CHECK ROUTES
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	})

	r.Get("/health/db", func(w http.ResponseWriter, r *http.Request) {
		if err := s.db.Ping(); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("db: not ok"))
			return
		}
		w.Write([]byte("db: ok"))
	})

	// ✅ APPROVAL REQUEST ROUTES
	r.Get("/requests", s.handleListRequests)
	r.Post("/requests", s.handleCreateRequest)
	r.Get("/requests/{id}", s.handleGetRequest)
	r.Post("/requests/{id}/approve", s.handleApproveRequest)
	r.Post("/requests/{id}/reject", s.handleRejectRequest)
	r.Get("/requests/{id}/logs", s.handleGetRequestLogs)

	s.router = r
}
