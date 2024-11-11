package api

import "github.com/go-chi/chi/v5"

func (a *API) addRoutes(mux *chi.Mux) {
	mux.Post("/api/v1/track", authorise(a.handlePostTrack))
	// mux.Get("/api/v1/count-by-event", authorise(a.handleGetEventCount))
}
