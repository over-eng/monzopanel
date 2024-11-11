package kafkatools

import (
	"context"
	"errors"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/rs/zerolog"
)

type TopicSpec struct {
	Topic             string `yaml:"topic"`
	NumPartitions     int    `yaml:"num_partitions"`
	ReplicationFactor int    `yaml:"replication_factor"`
}

func (t TopicSpec) ToKafkaTopicSpec() kafka.TopicSpecification {
	return kafka.TopicSpecification{
		Topic:             t.Topic,
		NumPartitions:     t.NumPartitions,
		ReplicationFactor: t.ReplicationFactor,
	}
}

// creates or updates topic to supplied spec
func EnsureTopic(ctx context.Context, client *kafka.AdminClient, topicSpec TopicSpec) error {
	log := zerolog.Ctx(ctx).With().Str("topic", topicSpec.Topic).Logger()
	res, err := client.DescribeTopics(ctx, kafka.NewTopicCollectionOfTopicNames([]string{topicSpec.Topic}))
	if err != nil {
		return err
	}

	description := res.TopicDescriptions[0]

	if description.Error.Code() == kafka.ErrUnknownTopicOrPart {
		log.Info().Msg("topic does not exist -> creating...")
		res, err := client.CreateTopics(
			ctx,
			[]kafka.TopicSpecification{topicSpec.ToKafkaTopicSpec()},
		)
		if err != nil {
			log.Err(err).Msg("failed to create topic")
			return err
		}

		if res[0].Error.Code() != kafka.ErrNoError {
			log.Err(err).Msg("failed to create topic")
			return errors.Join(errors.New("failed to create topic"), errors.New(res[0].Error.Error()))
		}
		log.Info().Msg("successfully created topic")
		return nil
	}

	if description.Error.Code() != kafka.ErrNoError {
		return errors.New(description.Error.Error())
	}

	partitionCount := len(description.Partitions)
	if partitionCount > topicSpec.NumPartitions {
		return errors.New("partition count cannot be decreased")
	}

	if len(description.Partitions) < topicSpec.NumPartitions {

		log.Info().
			Int("current_partition_count", partitionCount).
			Int("spec_partition_count", topicSpec.NumPartitions).
			Msg("partition count change detected, increasing...")
		res, err := client.CreatePartitions(ctx, []kafka.PartitionsSpecification{{
			Topic:      topicSpec.Topic,
			IncreaseTo: topicSpec.NumPartitions,
		}})
		if err != nil {
			log.Err(err).Msg("failed to increase partition count")
			return err
		}
		if res[0].Error.Code() != kafka.ErrNoError {
			return errors.New(res[0].Error.Error())
		}
	}

	return nil
}
