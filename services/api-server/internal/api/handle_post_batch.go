package api

import (
	"fmt"
	"io"
	"net/http"

	"github.com/gocql/gocql"
	"github.com/over-eng/monzopanel/protos/event"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (a *API) handlePostBatch(w http.ResponseWriter, r *http.Request) {

	body, err := io.ReadAll(r.Body)
	if err != nil {
		errorJSON(w, http.StatusBadRequest, "unable to read request body")
		return
	}

	batch := &event.EventBatch{}
	unmarshalOpts := protojson.UnmarshalOptions{
		DiscardUnknown: true,
		AllowPartial:   true,
	}
	if err := unmarshalOpts.Unmarshal(body, batch); err != nil {
		a.log.Err(err).Msg("failed to decode events")
		errorJSON(w, http.StatusBadRequest, "unable to decode events")
		return
	}

	events := batch.GetEvents()

	numEvents := len(events)
	if numEvents > a.config.Limits.MaxBatchSize {
		errorMessage := fmt.Sprintf("batch size exceeded, max %d", numEvents)
		errorJSON(w, http.StatusBadRequest, errorMessage)
		return
	}

	teamID := string(GetTeamIDFromRequest(r))
	createdAt := timestamppb.Now()
	for _, event := range events {
		if event.Id == "" {
			event.Id = gocql.TimeUUID().String()
		}
		event.TeamId = teamID
		event.CreatedAt = createdAt
	}

	res, err := a.eventwriter.WriteEvents(events)
	if err != nil {
		errorJSON(w, http.StatusInternalServerError, "unable to write events to kafka")
		return
	}

	if len(res.Fail) > 0 {
		encodeJSON(w, http.StatusInternalServerError, res)
		return
	}

	encodeJSON(w, http.StatusAccepted, res)
}
