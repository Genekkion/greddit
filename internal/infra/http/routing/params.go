package routing

import (
	"log/slog"

	servicesauth "greddit/internal/services/auth"
)

// RouterParams contains all the dependencies required by the router and sub
// routers.
type RouterParams struct {
	Logger *slog.Logger
	IsDev  bool

	AuthSer *servicesauth.Service
}

// Validate validates the dependencies of the router.
func (p RouterParams) Validate() error {
	if p.Logger == nil {
		return newInvalidRouterParamError("Logger")
	} else if p.AuthSer == nil {
		return newInvalidRouterParamError("AuthSer")
	}

	return nil
}
