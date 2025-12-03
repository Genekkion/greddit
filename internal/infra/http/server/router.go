package httpserver

import (
	"embed"
	"io/fs"
	"net/http"

	"greddit/internal/infra/http/api"

	"greddit/internal/infra/http/routing"
	"greddit/internal/infra/http/util"
)

// staticFs contains the static files to be served.
//
//go:embed static/*
var staticFs embed.FS

// NewHandler creates a new router.
func NewHandler(p routing.RouterParams) *http.ServeMux {
	mux := http.NewServeMux()

	if p.IsDev {
		sFs, err := fs.Sub(staticFs, "static")
		if err != nil {
			panic(err)
		}
		mux.Handle("/static/", http.StripPrefix("/static/",
			http.FileServer(http.FS(sFs))),
		)
	}

	httputil.AddSubRouters(mux, map[string]http.Handler{
		"/api": api.New(p),
	})

	mux.HandleFunc("/health", httputil.Methods(map[string]http.HandlerFunc{
		http.MethodGet: health,
	}))

	httputil.AddNotFoundHandler(mux)

	return mux
}

// health is the health check endpoint.
func health(w http.ResponseWriter, _ *http.Request) {
	httputil.WriteJson(w, http.StatusOK, map[string]any{
		"status": "healthy",
	})
}
