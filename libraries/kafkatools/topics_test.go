package kafkatools

import (
	"context"
	"testing"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/over-eng/monzopanel/libraries/testtools"
	"github.com/stretchr/testify/suite"
)

type topicSuite struct {
	suite.Suite
	kafka *testtools.KafkaSuite
}

func (suite *topicSuite) SetupSuite() {
	ctx := context.Background()

	kafkasuite, err := testtools.NewKafkaSuite(ctx)
	suite.Require().NoError(err)

	suite.kafka = kafkasuite
}

func (suite *topicSuite) TearDownSuite() {
	err := suite.kafka.TearDownSuite()
	suite.Require().NoError(err)
}

// The same container runs for all tests but we clear out the
// topics after each test.
func (suite *topicSuite) TearDownTest() {
	meta, err := suite.kafka.AdminClient.GetMetadata(nil, true, 10)
	suite.Require().NoError(err)

	topics := []string{}
	for _, topic := range meta.Topics {
		topics = append(topics, topic.Topic)
	}

	results, err := suite.kafka.AdminClient.DeleteTopics(context.Background(), topics)
	suite.Require().NoError(err)

	for _, r := range results {
		suite.Assert().Equal(r.Error.Code(), kafka.ErrNoError)
	}
}

func TestTopicSuite(t *testing.T) {
	suite.Run(t, new(topicSuite))
}

func (suite *topicSuite) TestEnsureTopicCreatesTopicIfNotExists() {
	ctx := context.Background()

	topic := TopicSpec{
		Topic:             "test-topic",
		NumPartitions:     3,
		ReplicationFactor: 1,
	}
	err := EnsureTopic(ctx, suite.kafka.AdminClient, topic)
	suite.Require().NoError(err)

	res, err := suite.kafka.AdminClient.DescribeTopics(
		ctx,
		kafka.NewTopicCollectionOfTopicNames([]string{"test-topic"}),
	)
	suite.Require().NoError(err)

	actual := res.TopicDescriptions[0]
	suite.Assert().Equal("test-topic", actual.Name)
	suite.Assert().Equal(topic.NumPartitions, len(actual.Partitions))
	suite.Assert().Equal(topic.ReplicationFactor, len(actual.Partitions[0].Replicas))
}

func (suite *topicSuite) TestEnsureTopicIsIdempotent() {
	ctx := context.Background()

	topic := TopicSpec{
		Topic:             "test-topic",
		NumPartitions:     3,
		ReplicationFactor: 1,
	}

	// create topic
	err := EnsureTopic(ctx, suite.kafka.AdminClient, topic)
	suite.Require().NoError(err)

	res, err := suite.kafka.AdminClient.DescribeTopics(
		ctx,
		kafka.NewTopicCollectionOfTopicNames([]string{"test-topic"}),
	)
	suite.Require().NoError(err)

	suite.Assert().Equal(kafka.ErrNoError, res.TopicDescriptions[0].Error.Code())

	actual := res.TopicDescriptions[0]
	suite.Assert().Equal("test-topic", actual.Name)
	suite.Assert().Equal(topic.NumPartitions, len(actual.Partitions))
	suite.Assert().Equal(topic.ReplicationFactor, len(actual.Partitions[0].Replicas))

	// run again to check repeated calls do not change state or error
	err = EnsureTopic(ctx, suite.kafka.AdminClient, topic)
	suite.Require().NoError(err)

	res, err = suite.kafka.AdminClient.DescribeTopics(
		ctx,
		kafka.NewTopicCollectionOfTopicNames([]string{"test-topic"}),
	)
	suite.Require().NoError(err)

	actual = res.TopicDescriptions[0]
	suite.Assert().Equal("test-topic", actual.Name)
	suite.Assert().Equal(topic.NumPartitions, len(actual.Partitions))
	suite.Assert().Equal(topic.ReplicationFactor, len(actual.Partitions[0].Replicas))
}

func (suite *topicSuite) TestEnsureTopicCanIncreasePartitions() {
	ctx := context.Background()

	topic := TopicSpec{
		Topic:             "test-topic",
		NumPartitions:     3,
		ReplicationFactor: 1,
	}
	// create topic
	err := EnsureTopic(ctx, suite.kafka.AdminClient, topic)
	suite.Require().NoError(err)

	res, err := suite.kafka.AdminClient.DescribeTopics(
		ctx,
		kafka.NewTopicCollectionOfTopicNames([]string{"test-topic"}),
	)
	suite.Require().NoError(err)

	actual := res.TopicDescriptions[0]
	suite.Assert().Equal("test-topic", actual.Name)
	suite.Assert().Equal(topic.NumPartitions, len(actual.Partitions))
	suite.Assert().Equal(topic.ReplicationFactor, len(actual.Partitions[0].Replicas))

	// increase partitions
	topic.NumPartitions = 6
	err = EnsureTopic(ctx, suite.kafka.AdminClient, topic)
	suite.Require().NoError(err)

	res, err = suite.kafka.AdminClient.DescribeTopics(
		ctx,
		kafka.NewTopicCollectionOfTopicNames([]string{"test-topic"}),
	)
	suite.Require().NoError(err)

	actual = res.TopicDescriptions[0]
	suite.Assert().Equal("test-topic", actual.Name)
	suite.Assert().Equal(topic.NumPartitions, len(actual.Partitions))
	suite.Assert().Equal(topic.ReplicationFactor, len(actual.Partitions[0].Replicas))
}
