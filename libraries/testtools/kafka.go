package testtools

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/testcontainers/testcontainers-go"
	kcontainer "github.com/testcontainers/testcontainers-go/modules/kafka"
)

type KafkaSuite struct {
	Container   *kcontainer.KafkaContainer
	Hosts       []string
	AdminClient *kafka.AdminClient
}

// Centralise this so we can update test images in one location.
func NewKafkaSuite(ctx context.Context) (*KafkaSuite, error) {

	container, err := kcontainer.Run(ctx, "confluentinc/confluent-local:7.7.1")
	if err != nil {
		return nil, err
	}

	c, err := container.Inspect(ctx)
	if err != nil {
		return nil, err
	}
	ports := c.NetworkSettings.Ports

	servers := []string{}
	for _, v := range ports {
		for _, i := range v {
			servers = append(servers, fmt.Sprintf("%s:%v", "127.0.0.1", i.HostPort))
		}
	}

	configMap := &kafka.ConfigMap{"bootstrap.servers": strings.Join(servers, ",")}
	admin, err := kafka.NewAdminClient(configMap)
	if err != nil {
		return nil, err
	}

	suite := &KafkaSuite{
		Container:   container,
		AdminClient: admin,
		Hosts:       servers,
	}
	return suite, nil

}

func (k *KafkaSuite) TearDownTest() error {
	meta, err := k.AdminClient.GetMetadata(nil, true, 10)
	if err != nil {
		return err
	}

	topics := []string{}
	for _, topic := range meta.Topics {
		topics = append(topics, topic.Topic)
	}

	results, err := k.AdminClient.DeleteTopics(context.Background(), topics)
	if err != nil {
		return err
	}

	for _, r := range results {
		if r.Error.Code() != kafka.ErrNoError {
			return errors.New(r.Error.Error())
		}
	}
	return nil
}

func (k *KafkaSuite) TearDownSuite() error {
	k.AdminClient.Close()
	return testcontainers.TerminateContainer(k.Container)
}

func (k *KafkaSuite) ConsumeTopic(ctx context.Context, topic string, numMessages int) ([]*kafka.Message, error) {
	config := &kafka.ConfigMap{
		"bootstrap.servers": strings.Join(k.Hosts, ","),
		"group.id":          "test-group",
		"auto.offset.reset": "earliest",
	}
	client, err := kafka.NewConsumer(config)
	if err != nil {
		return nil, err
	}
	defer client.Close()

	err = client.SubscribeTopics([]string{topic}, nil)
	if err != nil {
		return nil, err
	}

	messages := []*kafka.Message{}
	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
			ev, err := client.ReadMessage(100 * time.Second)
			if err != nil {
				return nil, err
			}
			messages = append(messages, ev)
		}

		if len(messages) >= numMessages {
			return messages, nil
		}
	}
}
