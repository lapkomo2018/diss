package server

import (
	"net/http"

	"lapkomo2018/diss/internal/server/handler"
)

type (
	Server struct {
		addr string

		serveMux http.Handler
	}
)

func NewServer(addr string) *Server {
	return &Server{
		addr:     addr,
		serveMux: http.NewServeMux(),
	}
}

func (s *Server) Init(p handler.Params) error {
	var err error
	if s.serveMux, err = handler.Init(p); err != nil {
		return err
	}

	return nil
}

func (s *Server) Server() *http.Server {
	return &http.Server{
		Addr:    s.addr,
		Handler: s.serveMux,
	}
}

func (s *Server) Start() error {
	return s.Server().ListenAndServe()
}
