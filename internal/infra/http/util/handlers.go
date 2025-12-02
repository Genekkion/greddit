package httputil

import (
	"net/http"
	"strings"
)

// allMethods is a list of all HTTP methods.
var allMethods = []string{
	http.MethodGet,
	http.MethodPut,
	http.MethodPost,
	http.MethodDelete,
	http.MethodHead,
	http.MethodPatch,
	http.MethodConnect,
	http.MethodOptions,
	http.MethodTrace,
}

// AddNotFoundHandler adds a handler for 404 errors which responds in JSON.
func AddNotFoundHandler(mux *http.ServeMux) {
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		m := r.Method
		_, curr := mux.Handler(r)

		var allowed []string

		for _, method := range allMethods {
			r.Method = method
			_, pattern := mux.Handler(r)
			if pattern != curr {
				allowed = append(allowed, method)
			}
		}

		r.Method = m

		if len(allowed) != 0 {
			w.Header().Set("allow", strings.Join(allowed, ", "))
			GenericMethodNotAllowed(w, r)
			return
		}

		GenericNotFound(w, r)
	})
}
