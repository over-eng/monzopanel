package eventwriter

import (
	"encoding/json"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/over-eng/monzopanel/protos/event"
)

type WriteEventsResult struct {
	Success []string `json:"success"`
	Fail    []string `json:"fail"`
	Invalid []string `json:"invalid"`
}

func (ew *EventWriter) WriteEvents(events []*event.Event) (WriteEventsResult, error) {

	result := WriteEventsResult{}
	delieveryChan := make(chan kafka.Event)

	for _, event := range events {

		err := event.ValidateQueueable()
		if err != nil {
			result.Invalid = append(result.Invalid, event.Id)
			continue
		}

		serialized, err := json.Marshal(event)
		if err != nil {
			ew.log.Error().Any("event", event).Msg("failed to serialize event")
			result.Fail = append(result.Invalid, event.Id)
			continue
		}

		message := &kafka.Message{
			TopicPartition: kafka.TopicPartition{
				Topic:     &ew.config.ProducerTopic.Topic,
				Partition: kafka.PartitionAny,
			},
			Value: serialized,
			Headers: []kafka.Header{
				{
					Key:   "attempts",
					Value: []byte("1"),
				},
			},
			Key: []byte(event.Id),
		}

		err = ew.producer.Produce(message, delieveryChan)
		if err != nil {
			ew.log.Error().Any("event", event).Msg("failed to write event to topic")
			result.Fail = append(result.Fail, event.Id)
			continue
		}
	}

	for {
		e := <-delieveryChan
		m := e.(*kafka.Message)
		if m.TopicPartition.Error != nil {
			result.Fail = append(result.Fail, string(m.Key))
		} else {
			result.Success = append(result.Success, string(m.Key))
		}
		if len(result.Success)+len(result.Fail) >= len(events) {
			break
		}
	}

	close(delieveryChan)

	return result, nil
}
