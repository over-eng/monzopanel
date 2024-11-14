package config

import (
	"os"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/over-eng/monzopanel/libraries/kafkatools"
	"github.com/over-eng/monzopanel/services/eventconsumer/internal/eventstore"
	"go.mau.fi/zeroconfig"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Metrics    Metrics           `yaml:"metrics"`
	Logging    zeroconfig.Config `yaml:"logging"`
	EventStore eventstore.Config `yaml:"event_store"`
	Broker     Broker            `yaml:"broker"`
}

type Metrics struct {
	Addr string `yaml:"addr"`
}

type Broker struct {
	EventConsumer      Consumer `yaml:"event_consumer"`
	RetryProducer      Producer `yaml:"retry_producer"`
	DeadLetterProducer Producer `yaml:"dead_letter_producer"`
}

type Consumer struct {
	Topic              string          `yaml:"topic"`
	ConfigMap          kafka.ConfigMap `yaml:"config_map"`
	AttemptsBeforeDead int             `yaml:"attempts_before_dead_letter"`
}

type Producer struct {
	TopicSpec kafkatools.TopicSpec `yaml:"topic_spec"`
	ConfigMap kafka.ConfigMap      `yaml:"config_map"`
}

func NewDefaultConfig() Config {
	return Config{
		Logging: zeroconfig.Config{
			Writers: []zeroconfig.WriterConfig{
				{
					Type: zeroconfig.WriterTypeStderr,
				},
			},
		},
	}
}

func (c *Config) LoadConfigFile(filename string) {
	data, err := os.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	err = yaml.Unmarshal(data, &c)
	if err != nil {
		panic(err)
	}
}
