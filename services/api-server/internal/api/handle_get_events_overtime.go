package api

import (
	"fmt"
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

	fromStr := r.URL.Query().Get("from")
	var from time.Time
	var err error
	if fromStr != "" {
		from, err = time.Parse(time.RFC3339, fromStr)
		if err != nil {
			errorJSON(w, http.StatusBadRequest, "invalid 'from' date format, expected RFC3339")
		}
	} else {
		from = time.Now().AddDate(0, 0, -7)
	}

	toStr := r.URL.Query().Get("to")
	var to time.Time
	if toStr != "" {
		to, err = time.Parse(time.RFC3339, toStr)
		if err != nil {
			errorJSON(w, http.StatusBadRequest, "invalid 'to' date format, expected RFC3339")
		}
	} else {
		to = time.Now()
	}

	periodStr := strings.ToUpper(r.URL.Query().Get("period"))
	if periodStr == "" {
		periodStr = "HOUR"
	}
	period, ok := event.TimePeriod_value[periodStr]
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
		To:     timestamppb.New(to),
		From:   timestamppb.New(from),
	}

	res, err := a.queryAPIClient.EventCountOvertime(r.Context(), &grpcReq)
	if err != nil {
		a.log.Err(err).Msg("failed to request event overtime count")
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
