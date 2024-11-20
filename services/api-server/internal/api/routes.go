package api

import "github.com/go-chi/chi/v5"

func (a *API) addRoutes(mux *chi.Mux) {
	mux.Post("/analytics/batch", authorise(a.handlePostBatch))
	// mux.Get("/api/v1/count-by-event", authorise(a.handleGetEventCount))
}
