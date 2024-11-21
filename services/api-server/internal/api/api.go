package api

import (
	"context"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/over-eng/monzopanel/services/api-server/internal/config"
	"github.com/over-eng/monzopanel/services/api-server/internal/eventwriter"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type API struct {
	log         zerolog.Logger
	config      config.Server
	server      *http.Server
	eventwriter *eventwriter.EventWriter
}

func New(cfg config.Server, eventwriter *eventwriter.EventWriter) *API {
	return &API{
		log:         log.With().Str("component", "api").Logger(),
		config:      cfg,
		eventwriter: eventwriter,
	}
}

func (a *API) Start() {
	mux := chi.NewRouter()
	mux.Use(a.corsMiddleware)
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
}

func (a *API) Stop(ctx context.Context) {
	err := a.server.Shutdown(ctx)
	if err != nil {
		a.log.Err(err).Msg("error shutting down server")
	}
}
