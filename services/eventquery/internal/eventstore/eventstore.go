package eventstore

import (
	"context"

	"github.com/gocql/gocql"
	"github.com/over-eng/monzopanel/libraries/cassandratools"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type Store struct {
	log     zerolog.Logger
	session *gocql.Session
}

type Config struct {
	Connection cassandratools.ConnectionConfig `yaml:"connection"`
	Keyspace   cassandratools.Keyspace         `yaml:"keyspace"`
}

func New(cfg Config) (*Store, error) {
	log := log.With().Str("component", "eventstore").Logger()
	ctx := log.WithContext(context.Background())
	log.Info().Msg("starting database session")

	session, err := cassandratools.NewSession(cfg.Connection).
		WithCreateKeyspace(cfg.Keyspace).
		WithUseKeyspace(cfg.Keyspace.Name).
		Start(ctx)
	if err != nil {
		log.Err(err).Msg("error starting database session")
		return nil, err
	}

	store := &Store{
		session: session,
		log:     log,
	}
	log.Info().Msg("database session created")
	return store, nil
}

func (s *Store) Close() {
	s.log.Info().Msg("shutting database session")
	s.session.Close()
}
