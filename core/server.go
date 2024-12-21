package core

import (
	"context"
	"errors"
	"net/http"
	"todo-ai/core/logger"
)

type Server interface {
	Start()
	Stop()
}

type server struct {
	*http.Server
}

func NewServer(addr string, handler http.Handler) Server {
	svr := &server{
		Server: &http.Server{
			Addr:    addr,
			Handler: handler,
		},
	}

	return svr
}

func (s *server) Start() {
	logger.Infof("HttpServer start with addr: %s", s.Addr)

	go func() {
		if err := s.ListenAndServe(); err != nil {
			if !errors.Is(err, http.ErrServerClosed) {
				logger.Errorf("HttpServer listen with err: %v", err)
			}
		}
	}()
}

func (s *server) Stop() {
	err := s.Shutdown(context.Background())
	if err != nil {
		logger.Errorf("HttpServer stop with err: %v", err)
	}
}
