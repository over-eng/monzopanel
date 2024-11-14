package consumer

import (
	"context"
	"encoding/json"
	"strconv"
	"sync"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/over-eng/monzopanel/libraries/kafkatools"
	"github.com/over-eng/monzopanel/libraries/models"
	"github.com/over-eng/monzopanel/services/eventconsumer/internal/config"
	"github.com/over-eng/monzopanel/services/eventconsumer/internal/eventstore"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type EventConsumer struct {
	log        zerolog.Logger
	eventstore *eventstore.Store
	config     config.Broker

	wg           sync.WaitGroup
	stopConsumer context.CancelFunc

	consumer   *kafka.Consumer
	retry      *kafka.Producer
	deadletter *kafka.Producer
}

func New(ctx context.Context, cfg config.Broker, store *eventstore.Store) (*EventConsumer, error) {
	log := log.With().Str("component", "event_consumer").Logger()

	err := ensureTopic(ctx, cfg.RetryProducer)
	if err != nil {
		log.Err(err).Msg("failed to ensure retry topic")
		return nil, err
	}

	err = ensureTopic(ctx, cfg.DeadLetterProducer)
	if err != nil {
		log.Err(err).Msg("failed to ensure dead letter topic")
		return nil, err
	}

	consumer := &EventConsumer{
		log:        log,
		eventstore: store,
		config:     cfg,
	}

	return consumer, nil
}

func (c *EventConsumer) Start() error {
	consumer, err := kafka.NewConsumer(&c.config.EventConsumer.ConfigMap)
	if err != nil {
		log.Err(err).Msg("failed to create event consumer")
		return err
	}
	c.consumer = consumer

	retry, err := kafka.NewProducer(&c.config.RetryProducer.ConfigMap)
	if err != nil {
		log.Err(err).Msg("failed to create retry producer")
		return err
	}
	c.retry = retry

	deadletter, err := kafka.NewProducer(&c.config.DeadLetterProducer.ConfigMap)
	if err != nil {
		log.Err(err).Msg("failed to create dead letter producer")
		return err
	}
	c.deadletter = deadletter

	ctx, cancel := context.WithCancel(context.Background())
	c.stopConsumer = cancel
	go c.consumeEvents(ctx)

	return nil
}

func (c *EventConsumer) Stop() {
	c.log.Info().Msg("stopping consumer")
	c.stopConsumer()
	c.wg.Wait()
	c.log.Info().Msg("consumer stopped")

	err := c.consumer.Close()
	if err != nil {
		c.log.Err(err).Msg("error shutting down consumer")
	}
	c.retry.Close()
	c.deadletter.Close()
}

func (c *EventConsumer) consumeEvents(ctx context.Context) {
	c.wg.Add(1)
	defer c.wg.Done()
	c.log.Info().Msg("starting consumer")

	topics := []string{c.config.EventConsumer.Topic, c.config.RetryProducer.TopicSpec.Topic}
	err := c.consumer.SubscribeTopics(topics, nil)
	if err != nil {
		c.log.Panic().Err(err).Msg("failed to subscribe to topics")
	}

	for {
		select {
		case <-ctx.Done():
			return
		default:
			m, err := c.consumer.ReadMessage(1 * time.Second)
			// log unexpected errors
			if err != nil && !err.(kafka.Error).IsTimeout() {
				c.log.Err(err).Msg("error consuming events")
				continue
			}
			if m == nil {
				c.log.Debug().Msg("no messages received within timeout")
				continue
			}
			c.log.Debug().Any("message", m).Msg("processing message")
			c.processMessage(ctx, m)
		}
	}

}

func (c *EventConsumer) processMessage(ctx context.Context, m *kafka.Message) {

	var failProducer *kafka.Producer
	for i, h := range m.Headers {
		if h.Key == "attempts" {
			attempts, err := strconv.Atoi(string(h.Value))
			if err != nil {
				c.log.Panic().Err(err).Msg("failed to convert attempt header to an int, this should never happen")
			}

			if attempts >= c.config.EventConsumer.AttemptsBeforeDead {
				failProducer = c.deadletter
			} else {
				failProducer = c.retry
			}

			// increment header now incase of retry
			m.Headers[i].Value = []byte(strconv.Itoa(attempts + 1))
		}
	}

	var event models.Event
	err := json.Unmarshal(m.Value, &event)
	if err != nil {
		c.log.Err(err).Msg("failed to unmarshal event")
		go c.handleFailedMessage(m, failProducer)
		return
	}

	err = c.eventstore.InsertEvent(ctx, &event)
	if err != nil {
		c.log.Err(err).Msg("failed to insert event into cassandra")
		go c.handleFailedMessage(m, failProducer)
		return
	}
}

func (c *EventConsumer) handleFailedMessage(m *kafka.Message, p *kafka.Producer) {
	c.log.Info().Any("event", m).Msg("retrying/dead-lettering event")
	delieveryChan := make(chan kafka.Event)
	err := p.Produce(m, delieveryChan)
	if err != nil {
		c.log.Err(err).Msg("failed to add event to retry/deadletter topic")
	}

	<-delieveryChan
	c.log.Info().Any("event", m).Msg("event delivered to retry/deadletter topic")

	close(delieveryChan)
}

func ensureTopic(ctx context.Context, cfg config.Producer) error {
	client, err := kafka.NewAdminClient(&cfg.ConfigMap)
	if err != nil {
		return err
	}
	defer client.Close()

	return kafkatools.EnsureTopic(ctx, client, cfg.TopicSpec)
}
