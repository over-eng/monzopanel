package api

import (
	"context"
	"net/http"
	"strings"
)

type CtxTeamID string

const teamIDContextKey CtxTeamID = "team_id"

func authorise(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		teamID := getTeamIdFromBearer(r)
		if teamID == "" {
			http.NotFound(w, r)
			return
		}

		// add teamID to context so it can be retrieved by handlers
		ctx := context.WithValue(r.Context(), teamIDContextKey, teamID)
		newReq := r.WithContext(ctx)

		next(w, newReq)
	})
}

func getTeamIdFromBearer(r *http.Request) string {
	var teamID string
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return teamID
	}

	token := strings.TrimPrefix(authHeader, "Bearer ")
	if token == "" {
		return teamID
	}

	// this is a major hack to get this going, but should be replaced by
	// teams stored in the db, and ideally and in memory ttl cache.
	if token == "the-super-secret-token" {
		teamID = "over-engineering.co.uk"
	}

	return teamID
}
