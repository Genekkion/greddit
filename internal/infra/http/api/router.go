package api

import (
	"net/http"

	v1 "greddit/internal/infra/http/api/v1"

	"greddit/internal/infra/http/routing"
	httputil "greddit/internal/infra/http/util"
)

// New creates a new router with the api routes.
func New(p routing.RouterParams) http.Handler {
	mux := http.NewServeMux()

	httputil.AddSubRouters(mux, map[string]http.Handler{
		"/v1": v1.ApiRoutes(p),
	})

	httputil.AddNotFoundHandler(mux)

	return mux
}
