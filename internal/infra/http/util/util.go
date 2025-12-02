package httputil

import (
	"encoding/json"
	"net/http"

	"greddit/internal/infra/log"
)

var logger = log.GetDefaultLogger()

// WriteJson writes a JSON response, checking for errors when writing.
func WriteJson(w http.ResponseWriter, status int, value any) {
	err := WriteJsonFull(w, status, value)
	if err != nil {
		logger.Error("Error writing json response",
			"status", status,
			"value", value,
			"error", err,
		)
	}
}

// WriteJsonFull writes a JSON response to the client.
func WriteJsonFull(w http.ResponseWriter, status int, value any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)

	err := json.NewEncoder(w).Encode(value)
	if err != nil {
		// This branch can probably be removed
		logger.Error("Error writing json response", "error", err)
	}
	return err
}
