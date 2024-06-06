package internalhttp

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"strconv"

	"github.com/Lanworm/OTUS_GO/hw12_13_14_15_calendar/internal/config"
	httpserver "github.com/Lanworm/OTUS_GO/hw12_13_14_15_calendar/internal/server"
)

type Server struct {
	logger Logger
	app    Application
	conf   config.ServerConf
	srv    *http.Server
	mux    *http.ServeMux
}

type Logger interface { // TODO
}

type Application interface { // TODO
}

func NewServer(
	logger Logger,
	app Application,
	conf config.ServerConf,
) *Server {
	if conf.Protocol == "" {
		conf.Protocol = "tcp4"
	}

	return &Server{
		logger: logger,
		app:    app,
		conf:   conf,
		mux:    http.NewServeMux(),
	}
}

func (s *Server) Start(_ context.Context) error {
	if s.srv != nil {
		return errors.New("server already started")
	}

	address := net.JoinHostPort(s.conf.Host, strconv.Itoa(s.conf.Port))
	s.srv = &http.Server{
		Addr:              address,
		Handler:           loggingMiddleware(s.mux),
		TLSConfig:         nil,
		ReadTimeout:       s.conf.Timeout,
		ReadHeaderTimeout: s.conf.Timeout,
		WriteTimeout:      s.conf.Timeout,
		IdleTimeout:       s.conf.Timeout,
		MaxHeaderBytes:    1 << 10,
	}

	registerRoutes(s)

	err := s.srv.ListenAndServe()
	if err != nil {
		return fmt.Errorf("listen and serve: %w", err)
	}

	return nil
}

func (s *Server) Stop(_ context.Context) error {
	return s.srv.Close()
}

func (s *Server) AddRoute(route string, handlerFunc http.HandlerFunc) {
	s.mux.HandleFunc(route, handlerFunc)
}

func registerRoutes(s *Server) {
	handler := new(httpserver.Handler)
	s.AddRoute("/hello", handler.Hello)
}
