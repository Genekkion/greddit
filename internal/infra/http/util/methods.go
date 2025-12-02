package httputil

import (
	"net/http"
	"slices"
)

// Methods map the specified allMethods to the specified handler.
// Panics if the allMethods param is empty (due to developer error).
func Methods(handlerMap map[string]http.HandlerFunc) http.HandlerFunc {
	if len(handlerMap) == 0 {
		panic("Requires at least 1 allowed method")
	}

	for k, v := range handlerMap {
		if !slices.Contains(allMethods, k) {
			panic("Unknown method: " + k)
		} else if v == nil {
			panic("Handler for method " + k + " is nil")
		}
	}

	return func(w http.ResponseWriter, r *http.Request) {
		handler, ok := handlerMap[r.Method]
		if !ok {
			GenericMethodNotAllowed(w, r)
			return
		}

		handler(w, r)
	}
}
