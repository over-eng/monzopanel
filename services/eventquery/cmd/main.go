package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"runtime/pprof"
	"time"

	"github.com/over-eng/monzopanel/services/eventquery/internal/config"
	"github.com/over-eng/monzopanel/services/eventquery/internal/eventstore"
	"github.com/over-eng/monzopanel/services/eventquery/internal/server"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func run(ctx context.Context, cfg config.Config) error {
	logger, err := cfg.Logging.Compile()
	if err != nil {
		return err
	}
	log.Logger = *logger
	zerolog.DefaultContextLogger = &log.Logger

	store, err := eventstore.New(cfg.EventStore)
	if err != nil {
		log.Err(err).Msg("failed to create event store")
		return err
	}
	server := server.NewServer(cfg.Server, store)
	server.Start()

	ctx, _ = signal.NotifyContext(ctx, os.Interrupt)
	<-ctx.Done()

	log.Info().Msg("stopping")

	// detect shutdown deadlocks and print goroutine dump
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func(ctx context.Context) {
		select {
		case <-ctx.Done():
			return
		case <-time.After(time.Minute):
			log.Error().Msg("detected a deadlock on shutdown, dumping goroutines:")
			pprof.Lookup("goroutine").WriteTo(os.Stderr, 2)
			os.Exit(1)
		}
	}(ctx)

	server.Stop()
	store.Close()

	log.Info().Msg("Exit 0")
	return nil
}

func main() {
	ctx := context.Background()

	configPath := flag.String("config", "config.yaml", "Config file")
	flag.Parse()

	cfg := config.NewDefaultConfig()
	cfg.LoadConfigFile(*configPath)

	if err := run(ctx, cfg); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}
