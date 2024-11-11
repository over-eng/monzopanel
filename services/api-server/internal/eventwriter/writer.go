package eventwriter

import (
	"context"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/over-eng/monzopanel/libraries/kafkatools"
	"github.com/over-eng/monzopanel/services/api-server/internal/config"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type EventWriter struct {
	log      zerolog.Logger
	producer *kafka.Producer
	config   config.Kafka
}

func New(ctx context.Context, cfg config.Kafka) (*EventWriter, error) {
	log := log.With().Str("component", "event_writer").Logger()

	log.Info().Any("config", cfg.ConfigMap).Msg("starting kafka")

	producer, err := kafka.NewProducer(&cfg.ConfigMap)
	if err != nil {
		log.Err(err).Msg("failed to create producer")
		return nil, err
	}

	writer := &EventWriter{
		log:      log,
		producer: producer,
		config:   cfg,
	}

	adminClient, err := kafka.NewAdminClientFromProducer(producer)
	if err != nil {
		log.Err(err).Msg("failed create kafka admin client")
		return nil, err
	}
	defer adminClient.Close()

	err = kafkatools.EnsureTopic(ctx, adminClient, cfg.ProducerTopic)
	if err != nil {
		return nil, err
	}

	return writer, nil
}

func (ew *EventWriter) Close(timeout time.Duration) {
	ew.log.Info().Msg("flushing producer")
	unflushedCount := ew.producer.Flush(int(timeout.Milliseconds()))
	ew.log.Info().Int("unflushed_count", unflushedCount).Msg("finished flush")
	ew.producer.Close()
	ew.log.Info().Msg("producer closed")
}
