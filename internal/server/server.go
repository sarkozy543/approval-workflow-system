package server

import (
	"database/sql"
	"net/http"

	"github.com/sarkozy543/approval-workflow-system/internal/approval"
)

type Server struct {
	router http.Handler
	db     *sql.DB

	approvalStore *approval.Store
}

// NewServer: dışarıdan DB alır ve store'ları hazırlar
func NewServer(db *sql.DB) *Server {
	s := &Server{
		db:            db,
		approvalStore: approval.NewStore(db),
	}
	s.routes()
	return s
}

func (s *Server) Router() http.Handler {
	return s.router
}
