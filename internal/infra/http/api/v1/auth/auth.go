package httpapiauth

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strings"

	httputil "greddit/internal/infra/http/util"

	"greddit/internal/infra/http/routing"
	servicesauth "greddit/internal/services/auth"
)

type AuthRouter struct {
	logger *slog.Logger
	ser    servicesauth.Service
}

func AuthRoutes(p routing.RouterParams) (mux *http.ServeMux) {
	mux = http.NewServeMux()

	rtr := AuthRouter{
		ser:    *p.AuthSer,
		logger: p.Logger,
	}

	mux.HandleFunc("/login", httputil.Methods(map[string]http.HandlerFunc{
		http.MethodPost: rtr.login,
	}))

	mux.HandleFunc("/check", httputil.Methods(map[string]http.HandlerFunc{
		http.MethodGet: rtr.checkAuth,
	}))

	return mux
}

func (rtr AuthRouter) login(w http.ResponseWriter, r *http.Request) {
	var reqBody struct {
		Username string `json:"username"`
	}
	defer r.Body.Close()

	err := json.NewDecoder(r.Body).Decode(&reqBody)
	if err != nil {
		rtr.logger.ErrorContext(r.Context(), "Error decoding request body",
			"error", err,
		)
		httputil.GenericBadRequest(w, r)
		return
	}

	signed, err := rtr.ser.Login(r.Context(), reqBody.Username)
	if err != nil {
		rtr.logger.ErrorContext(r.Context(), "Error signing JWT token",
			"error", err,
		)
		httputil.GenericInternalServerError(w, r)
		return
	}

	httputil.WriteJson(w, http.StatusOK, map[string]any{
		"token": string(signed),
	})
}

func (rtr AuthRouter) checkAuth(w http.ResponseWriter, r *http.Request) {
	value := r.Header.Get("Authorization")
	if value == "" {
		rtr.logger.ErrorContext(r.Context(), "No token provided")
		httputil.GenericUnauthorized(w, r)
		return
	}
	value = strings.TrimSpace(value)
	value = strings.TrimPrefix(value, "Bearer ")
	value = strings.TrimSpace(value)

	rtr.logger.DebugContext(r.Context(), "Token received",
		"value", value,
	)

	token, err := rtr.ser.ExtractClaims(r.Context(), []byte(value))
	if err != nil {
		rtr.logger.ErrorContext(r.Context(), "Error extracting claims from token",
			"error", err,
		)
		httputil.GenericUnauthorized(w, r)
		return
	}
	rtr.logger.DebugContext(r.Context(), "Token extracted",
		"claims", token,
	)

	httputil.WriteJson(w, http.StatusOK, map[string]any{
		"message": "valid token",
	})
}
