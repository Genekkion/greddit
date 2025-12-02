package server

import (
	"log/slog"
	"net/http"

	"greddit/internal/infra/log"
)

// TracerMiddleware adds a trace ID to the request context.
func TracerMiddleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r = r.WithContext(log.CtxWithTraceID(r.Context()))

		handler.ServeHTTP(w, r)
	})
}

// LoggerMiddleware logs all incoming requests.
func LoggerMiddleware(handler http.Handler, logger *slog.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.DebugContext(r.Context(), "HTTP request",
			"method", r.Method,
			"host", r.Host,
			"url", r.URL.RequestURI(),
		)

		handler.ServeHTTP(w, r)
	})
}

// corsMiddleware adds CORS headers to the response.
func corsMiddleware(handler http.Handler, allowedOrigins string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", allowedOrigins)
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS, PATCH")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		handler.ServeHTTP(w, r)
	})
}
