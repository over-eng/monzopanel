package api

import (
	"net/http"

	"github.com/gogo/protobuf/jsonpb"
	"github.com/over-eng/monzopanel/protos/event"
)

func (a *API) handleGetEvents(w http.ResponseWriter, r *http.Request) {
	var req event.ListEventsByDistinctIDRequest
	unmarshaler := &jsonpb.Unmarshaler{}
	err := unmarshaler.Unmarshal(r.Body, &req)
	if err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	// we're only allowed to view our own team so overwrite the request
	// with the bearer's team id.
	req.TeamId = GetTeamIDFromRequest(r)

	response, err := a.queryAPIClient.ListEventsByDistinctID(r.Context(), &req)
	if err != nil {
		a.log.Error().Err(err).Msg("failed to list events")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	encodeJSON(w, http.StatusOK, response)
}
