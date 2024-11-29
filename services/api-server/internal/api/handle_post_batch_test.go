package api_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/over-eng/monzopanel/protos/event"
	"github.com/over-eng/monzopanel/services/api-server/internal/eventwriter"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (suite *testAPISuite) TestHandlePostTrack() {

	distinctID := suite.T().Name()

	events := &event.EventBatch{
		Events: []*event.Event{{
			Event:           "test-event",
			DistinctId:      distinctID,
			ClientTimestamp: timestamppb.Now(),
		}},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	messageCh := make(chan []*kafka.Message)
	go func() {
		messages, err := suite.kafka.ConsumeTopic(ctx, "test-topic", 1)
		suite.Assert().NoError(err)
		messageCh <- messages
	}()

	reqBody, err := protojson.Marshal(events)
	suite.Assert().NoError(err)

	request, err := http.NewRequest(http.MethodPost, "http://127.0.0.1:5999/analytics/batch", bytes.NewBuffer(reqBody))
	request.Header.Add("Authorization", "Bearer the-super-secret-token")
	suite.Assert().NoError(err)

	client := http.Client{}
	response, err := client.Do(request)
	suite.Assert().NoError(err)

	suite.Assert().Equal(http.StatusAccepted, response.StatusCode)

	var result eventwriter.WriteEventsResult
	err = json.NewDecoder(response.Body).Decode(&result)
	suite.Assert().NoError(err)

	suite.Assert().Equal(1, len(result.Success))
	suite.Assert().Equal(0, len(result.Invalid))
	suite.Assert().Equal(0, len(result.Fail))

	messages := <-messageCh
	suite.Assert().Equal(len(messages), 1)
}
