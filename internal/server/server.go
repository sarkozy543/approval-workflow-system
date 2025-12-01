package server

import "net/http"

type Server struct {
	router http.Handler
}

func NewServer() *Server {
	s := &Server{}
	s.routes()
	return s
}

func (s *Server) Router() http.Handler {
	return s.router
}
