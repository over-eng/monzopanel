package metrics

import (
	"net/http"

	"github.com/over-eng/monzopanel/services/eventconsumer/internal/config"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var (
	EventsInserted = promauto.NewCounter(prometheus.CounterOpts{
		Name: "event_insert_count",
		Help: "Counter for events inserted into cassandra",
	})
	EventInsertLatency = promauto.NewHistogram(prometheus.HistogramOpts{
		Name:    "event_insert_latency_seconds",
		Help:    "Latency for inserting an event into cassandra",
		Buckets: prometheus.LinearBuckets(0, 0.1, 51),
	})
	EventsRetried = promauto.NewCounter(prometheus.CounterOpts{
		Name: "event_retry_count",
		Help: "Counter for events added to retry topic",
	})
	EventsDeadLettered = promauto.NewCounter(prometheus.CounterOpts{
		Name: "event_dead_letter_count",
		Help: "Counter for events added to dead letter topic",
	})
)

type Server struct {
	log    zerolog.Logger
	server *http.Server
}

func NewServer(cfg config.Metrics) *Server {
	logger := log.With().
		Str("component", "metrics").
		Logger()

	mux := http.NewServeMux()
	mux.Handle("/", promhttp.Handler())
	mux.Handle("/metrics", promhttp.Handler())
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	return &Server{
		log:    logger,
		server: &http.Server{Addr: cfg.Addr, Handler: mux},
	}
}

func (mh *Server) Start() {
	mh.log.Info().Msgf("Starting metrics HTTP server at: %s", mh.server.Addr)
	go func() {
		err := mh.server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			mh.log.Fatal().Err(err).Msg("Error in metrics listener")
		}
	}()
}

func (mh *Server) Stop() {
	mh.log.Info().Msg("Stopping metrics HTTP server")
	err := mh.server.Close()
	if err != nil {
		mh.log.Err(err).Msg("Error closing metrics listener")
	}
}
