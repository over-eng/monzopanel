package config

import (
	"os"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/over-eng/monzopanel/libraries/kafkatools"
	"go.mau.fi/zeroconfig"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Server  Server            `yaml:"server"`
	Metrics Metrics           `yaml:"metrics"`
	Logging zeroconfig.Config `yaml:"logging"`
	Kafka   Kafka             `yaml:"kafka"`
}

type Server struct {
	Addr           string   `yaml:"addr"`
	Limits         Limits   `yaml:"limits"`
	AllowedOrigins []string `yaml:"allowed_origins"`
}

type Limits struct {
	MaxBatchSize int `yaml:"max_batch_size"`
}

type Metrics struct {
	Addr string `yaml:"addr"`
}

type Kafka struct {
	ConfigMap     kafka.ConfigMap      `yaml:"config_map"`
	ProducerTopic kafkatools.TopicSpec `yaml:"producer_topic"`
}

type Topic struct {
	Topic             string `yaml:"topic"`
	NumPartitions     int    `yaml:"numPartitions"`
	ReplicationFactor int    `yaml:"replicationFactor"`
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
