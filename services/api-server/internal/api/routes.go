package api

import "github.com/go-chi/chi/v5"

func (a *API) addRoutes(mux *chi.Mux) {
	mux.Post("/analytics/batch", authorise(a.handlePostBatch))

	mux.Get("/analytics/events", authorise(a.handleGetEvents))
	mux.Get("/analytics/stats/eventsovertime", authorise(a.handleGetEventsOvertime))
}
