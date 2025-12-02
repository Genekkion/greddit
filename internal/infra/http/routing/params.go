package routing

import (
	"log/slog"
)

// RouterParams contains all the dependencies required by the router and sub
// routers.
type RouterParams struct {
	Logger *slog.Logger

	IsDev bool
}

// Validate validates the dependencies of the router.
func (p RouterParams) Validate() error {
	return nil
}
