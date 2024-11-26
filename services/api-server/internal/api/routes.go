package api

import "github.com/go-chi/chi/v5"

func (a *API) addRoutes(mux *chi.Mux) {
	mux.Post("/analytics/batch", authorise(a.handlePostBatch))

	mux.Get("/analytics/{distinctId}/events", authorise(a.handleGetEvents))
	mux.Get("/analytics/{distinctId}/eventsovertime", authorise(a.handleGetEventsOvertime))
}
