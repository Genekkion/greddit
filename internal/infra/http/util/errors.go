package httputil

import "net/http"

// RespError writes an error response to the http.ResponseWriter.
func RespError(w http.ResponseWriter, _ *http.Request, statusCode int, reason string, logArgs ...any) {
	logArgs = append([]any{
		"status", statusCode,
		"reason", reason,
	}, logArgs...)
	logger.Error("HTTP Response error", logArgs...)

	WriteJson(w, statusCode, map[string]any{
		"error": reason,
	})
}

// GenericNotFound writes a generic 404 error response.
func GenericNotFound(w http.ResponseWriter, r *http.Request) {
	RespError(w, r, http.StatusNotFound, "not found")
}

// GenericMethodNotAllowed writes a generic 405 error response.
func GenericMethodNotAllowed(w http.ResponseWriter, r *http.Request) {
	RespError(w, r, http.StatusMethodNotAllowed, "method not allowed")
}

// GenericBadRequest writes a generic 400 error response.
func GenericBadRequest(w http.ResponseWriter, r *http.Request) {
	RespError(w, r, http.StatusBadRequest, "bad request")
}

// GenericInternalServerError writes a generic 500 error response.
func GenericInternalServerError(w http.ResponseWriter, r *http.Request) {
	RespError(w, r, http.StatusInternalServerError, "internal server error")
}

// GenericUnauthorized writes a generic 401 error response.
func GenericUnauthorized(w http.ResponseWriter, r *http.Request) {
	RespError(w, r, http.StatusUnauthorized, "unauthorized")
}
