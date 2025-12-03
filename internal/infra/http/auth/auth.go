package auth

import (
	"context"
	"net/http"

	httputil "greddit/internal/infra/http/util"

	servicesauth "greddit/internal/services/auth"
)

type CtxAuthKey int

const (
	CtxAuth CtxAuthKey = iota
)

// AuthMiddleware extracts the claims from the request and adds them to the context.
func AuthMiddleware(handler http.Handler, ser servicesauth.Service) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		claims, err := ser.ExtractClaims(r.Context(), []byte(r.Header.Get("Authorization")))
		if err != nil {
			httputil.RespError(w, r, http.StatusUnauthorized, "Invalid token")
			return
		}

		r = r.WithContext(context.WithValue(r.Context(), CtxAuth, *claims))

		handler.ServeHTTP(w, r)
	})
}

// GetClaims returns the claims from the context. Warning: to be used only
// if the above middleware has been applied to the request.
func GetClaims(r *http.Request) servicesauth.TokenClaims {
	v := r.Context().Value(CtxAuth)
	vv, ok := v.(servicesauth.TokenClaims)
	if !ok {
		panic("invalid claims type")
	}
	return vv
}
