package internalgrpc

import (
	"context"
	"errors"
	"fmt"
	"net"

	"github.com/Lanworm/OTUS_GO/hw12_13_14_15_calendar/internal/config"
	"github.com/Lanworm/OTUS_GO/hw12_13_14_15_calendar/internal/logger"
	grpc_calendar "github.com/Lanworm/OTUS_GO/hw12_13_14_15_calendar/pb/grpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type Server struct {
	handler grpc_calendar.CalendarServer
	conf    config.ServerGRPCConf
	logger  *logger.Logger
	server  *grpc.Server
}

func NewRPCServer(
	handler grpc_calendar.CalendarServer,
	logger *logger.Logger,
	conf config.ServerGRPCConf,
) *Server {
	if conf.Protocol == "" {
		conf.Protocol = "tcp4"
	}

	return &Server{
		handler: handler,
		conf:    conf,
		logger:  logger,
	}
}

func (s *Server) Start(_ context.Context) error {
	if s.server != nil {
		return errors.New("grpc server already started")
	}

	lsn, err := net.Listen(s.conf.Protocol, s.conf.GetFullAddress())
	if err != nil {
		return fmt.Errorf("grpc listen at {%s} : %w", s.conf.GetFullAddress(), err)
	}

	s.server = grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			LoggingRequest(s.logger),
		),
	)

	grpc_calendar.RegisterCalendarServer(s.server, s.handler)
	reflection.Register(s.server)

	if err := s.server.Serve(lsn); err != nil {
		return fmt.Errorf("serve grpc: %w", err)
	}

	return nil
}

func (s *Server) Stop(_ context.Context) error {
	s.server.GracefulStop()

	return nil
}
