package api

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func encodeJSON[T any](w http.ResponseWriter, status int, v T) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		return fmt.Errorf("encode json: %w", err)
	}
	return nil
}

func decodeJSON[T any](r *http.Request) (T, error) {
	var v T
	if err := json.NewDecoder(r.Body).Decode(&v); err != nil {
		return v, fmt.Errorf("decode json: %w", err)
	}
	return v, nil
}

func errorJSON(w http.ResponseWriter, status int, message string) error {
	json := struct {
		Error string `json:"error"`
	}{
		Error: message,
	}
	return encodeJSON(w, status, json)
}

func GetTeamIDFromRequest(r *http.Request) string {
	return r.Context().Value(teamIDContextKey).(string)
}
