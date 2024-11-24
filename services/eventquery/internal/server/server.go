package server

import (
	"net"

	"github.com/over-eng/monzopanel/protos/event"
	"github.com/over-eng/monzopanel/services/eventquery/internal/eventstore"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
)

type Server struct {
	event.UnimplementedQueryAPIServer
	log          zerolog.Logger
	server       *grpc.Server
	healthserver *health.Server

	config     Config
	eventstore *eventstore.Store
}

type Config struct {
	Addr              string `yaml:"addr"`
	ReflectionEnabled bool   `yaml:"enable_reflection"`
}

func NewServer(cfg Config, store *eventstore.Store) *Server {
	log := log.With().Str("component", "server").Logger()

	s := grpc.NewServer()

	return &Server{
		log:        log,
		config:     cfg,
		server:     s,
		eventstore: store,
	}
}

func (s *Server) Start() error {
	s.log.Info().Msg("starting server")
	lis, err := net.Listen("tcp", s.config.Addr)
	if err != nil {
		s.log.Err(err).
			Str("port", s.config.Addr).
			Msg("failed to listen on port, server not started")
		return err
	}

	s.healthserver = health.NewServer()
	grpc_health_v1.RegisterHealthServer(s.server, s.healthserver)

	event.RegisterQueryAPIServer(s.server, s)

	if s.config.ReflectionEnabled {
		s.log.Info().Msg("enabling reflection")
		reflection.Register(s.server)
	}

	go func() {
		err = s.server.Serve(lis)
		if err != nil {
			s.log.Fatal().Err(err).Msg("failed to serve")
		}
	}()

	s.log.Info().Msgf("server listening on: %s", s.config.Addr)
	s.healthserver.SetServingStatus("eventquery", grpc_health_v1.HealthCheckResponse_SERVICE_UNKNOWN)
	return nil
}

func (s *Server) Stop() {
	s.log.Info().Msg("shutting down server")
	s.healthserver.SetServingStatus("eventquery", grpc_health_v1.HealthCheckResponse_NOT_SERVING)
	s.server.GracefulStop()
	s.healthserver.Shutdown()
}
