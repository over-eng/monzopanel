package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/over-eng/monzopanel/protos/event"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type EventCountOvertimeRequest struct {
	To     time.Time `yaml:"to"`
	From   time.Time `yaml:"from"`
	Period string    `yaml:"period"`
}

type EventCountOvertimeResponse struct {
	Buckets []Bucket `json:"buckets"`
}

type Bucket struct {
	Timestamp time.Time `json:"timestamp"`
	Count     int64     `json:"count"`
}

func (a *API) handleGetEventsOvertime(w http.ResponseWriter, r *http.Request) {

	rawBody, err := io.ReadAll(r.Body)
	if err != nil {
		errorJSON(w, http.StatusBadRequest, "cannot read request json")
		return
	}

	reqBody := EventCountOvertimeRequest{}
	err = json.Unmarshal(rawBody, &reqBody)
	if err != nil {
		errorJSON(w, http.StatusBadRequest, "invalid JSON format")
		return
	}

	period, ok := event.TimePeriod_value[strings.ToUpper(reqBody.Period)]
	if !ok {
		validPeriods := []string{}
		for period := range event.TimePeriod_value {
			validPeriods = append(validPeriods, strings.ToLower(period))
		}
		errStr := fmt.Sprintf("period needs to be: %v", validPeriods)
		errorJSON(w, http.StatusBadRequest, errStr)
		return
	}

	grpcReq := event.EventCountOvertimeRequest{
		TeamId: GetTeamIDFromRequest(r),
		Period: event.TimePeriod(period),
		To:     timestamppb.New(reqBody.To),
		From:   timestamppb.New(reqBody.From),
	}

	res, err := a.queryAPIClient.EventCountOvertime(r.Context(), &grpcReq)
	if err != nil {
		errorJSON(w, http.StatusInternalServerError, "failed to request event overtime count")
		return
	}

	buckets := []Bucket{}
	for _, bucket := range res.GetBuckets() {
		buckets = append(buckets, Bucket{Count: bucket.GetCount(), Timestamp: bucket.GetTimestamp().AsTime()})
	}

	handlerResponse := EventCountOvertimeResponse{Buckets: buckets}
	encodeJSON(w, http.StatusOK, handlerResponse)
}
