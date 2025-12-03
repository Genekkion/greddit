package v1

import (
	"net/http"

	httpapiauth "greddit/internal/infra/http/api/v1/auth"
	httputil "greddit/internal/infra/http/util"

	"greddit/internal/infra/http/routing"
)

// ApiRoutes creates a new router with the api routes.
func ApiRoutes(p routing.RouterParams) (mux *http.ServeMux) {
	mux = http.NewServeMux()

	httputil.AddSubRouters(mux, map[string]http.Handler{
		"/auth": httpapiauth.AuthRoutes(p),
	})

	return mux
}
