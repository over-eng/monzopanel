package api

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/over-eng/monzopanel/protos/event"
)

func (a *API) handleGetEvents(w http.ResponseWriter, r *http.Request) {

	pageSize, err := strconv.Atoi(r.URL.Query().Get("page_size"))
	if err != nil {
		errorJSON(w, http.StatusBadRequest, "unable to parse page_size")
		return
	}

	if pageSize <= 0 {
		pageSize = 10
	}

	req := event.ListEventsByDistinctIDRequest{
		DistinctId:      chi.URLParam(r, "distinctId"),
		PageSize:        int32(pageSize),
		PaginationToken: r.URL.Query().Get("pagination_token"),
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
