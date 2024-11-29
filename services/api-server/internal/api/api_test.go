package api_test

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/over-eng/monzopanel/libraries/kafkatools"
	"github.com/over-eng/monzopanel/libraries/testtools"
	"github.com/over-eng/monzopanel/services/api-server/internal/api"
	"github.com/over-eng/monzopanel/services/api-server/internal/config"
	"github.com/over-eng/monzopanel/services/api-server/internal/eventwriter"
	"github.com/stretchr/testify/suite"
)

type testAPISuite struct {
	suite.Suite
	kafka  *testtools.KafkaSuite
	writer *eventwriter.EventWriter
	api    *api.API
}

func (suite *testAPISuite) SetupSuite() {
	kafkasuite, err := testtools.NewKafkaSuite(context.Background())
	suite.Require().NoError(err)
	suite.kafka = kafkasuite
}

func (suite *testAPISuite) SetupTest() {
	ctx := context.Background()
	writer, err := eventwriter.New(ctx, config.Kafka{
		ConfigMap: kafka.ConfigMap{
			"bootstrap.servers": strings.Join(suite.kafka.Hosts, ","),
			"acks":              "1",
		},
		ProducerTopic: kafkatools.TopicSpec{
			Topic:             "test-topic",
			NumPartitions:     1,
			ReplicationFactor: 1,
		},
	})
	suite.Require().NoError(err)
	suite.writer = writer

	cfg := config.Server{
		Addr:   ":5999",
		Limits: config.Limits{MaxBatchSize: 10},
	}

	suite.api = api.New(cfg, writer)
	suite.api.Start()
}

func (suite *testAPISuite) TearDownTest() {
	suite.writer.Close(2 * time.Second)
	suite.api.Stop(context.Background())
	err := suite.kafka.TearDownTest()
	suite.Require().NoError(err)
}

func (suite *testAPISuite) TearDownSuite() {
	err := suite.kafka.TearDownSuite()
	suite.Assert().NoError(err)
}

func TestAPITestSuite(t *testing.T) {
	suite.Run(t, new(testAPISuite))
}
