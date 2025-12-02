package httputil

import "net/http"

// AddSubRouters adds sub routers to the main router.
func AddSubRouters(m *http.ServeMux, subRouters map[string]http.Handler) {
	for k, v := range subRouters {
		m.Handle(k+"/", http.StripPrefix(k, v))
	}
}
