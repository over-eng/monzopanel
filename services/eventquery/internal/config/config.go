package config

import (
	"os"

	"github.com/over-eng/monzopanel/services/eventquery/internal/eventstore"
	"github.com/over-eng/monzopanel/services/eventquery/internal/server"
	"go.mau.fi/zeroconfig"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Metrics    Metrics           `yaml:"metrics"`
	Logging    zeroconfig.Config `yaml:"logging"`
	EventStore eventstore.Config `yaml:"event_store"`
	Server     server.Config     `yaml:"server"`
}

type Metrics struct {
	Addr string `yaml:"addr"`
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
