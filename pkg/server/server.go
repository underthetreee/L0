package server

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/underthetreee/L0/config"
)

type Server struct {
	srv *http.Server
}

func NewServer(cfg *config.Config, handler http.Handler) *Server {
	return &Server{
		srv: &http.Server{
			Addr:         ":" + cfg.HTTP.Port,
			Handler:      handler,
			ReadTimeout:  cfg.HTTP.ReadTimeout,
			WriteTimeout: cfg.HTTP.WriteTimeout,
		},
	}
}

func (s *Server) Run() error {
	log.Printf("listening on %s", s.srv.Addr)
	if err := s.srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("run http server: %w", err)
	}
	return nil
}

func (s *Server) Shutdown(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	log.Printf("shutting down server...")
	if err := s.srv.Shutdown(ctx); err != nil {
		return fmt.Errorf("shutdown http server: %w", err)
	}
	return nil
}
