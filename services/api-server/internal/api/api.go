package api

import (
	"context"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/over-eng/monzopanel/protos/event"
	"github.com/over-eng/monzopanel/services/api-server/internal/config"
	"github.com/over-eng/monzopanel/services/api-server/internal/eventwriter"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type API struct {
	log         zerolog.Logger
	config      config.Server
	server      *http.Server
	eventwriter *eventwriter.EventWriter

	queryAPIConn   *grpc.ClientConn
	queryAPIClient event.QueryAPIClient
}

func New(cfg config.Server, eventwriter *eventwriter.EventWriter) *API {
	return &API{
		log:         log.With().Str("component", "api").Logger(),
		config:      cfg,
		eventwriter: eventwriter,
	}
}

func (a *API) Start() error {
	a.log.Info().Msg("starting api server")

	conn, err := grpc.NewClient(
		a.config.QueryAPI.Host,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return err
	}
	a.queryAPIConn = conn
	a.queryAPIClient = event.NewQueryAPIClient(a.queryAPIConn)

	mux := chi.NewRouter()
	mux.Use(cors.Handler(cors.Options{
		AllowedOrigins:   a.config.AllowedOrigins,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))
	mux.Use(middleware.Logger)
	mux.Use(middleware.Heartbeat("/ping"))

	a.addRoutes(mux)

	a.server = &http.Server{
		Addr:    a.config.Addr,
		Handler: mux,
	}
	a.log.Info().Msgf("starting HTTP server at: %s", a.server.Addr)
	go func() {
		err := a.server.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			a.log.Fatal().Err(err).Msg("error while listening")
			return
		}
		a.log.Info().Msg("listener stopped")
	}()
	return nil
}

func (a *API) Stop(ctx context.Context) {
	err := a.queryAPIConn.Close()
	if err != nil {
		a.log.Err(err).Msg("error closing query api gRPC connection")
	}

	err = a.server.Shutdown(ctx)
	if err != nil {
		a.log.Err(err).Msg("error shutting down server")
	}
}
