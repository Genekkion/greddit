package routing

// InvalidRouterParamError represents an error when the router params are invalid.
type InvalidRouterParamError struct {
	ParamName string
}

// newInvalidRouterParamError creates a new InvalidRouterParamError.
func newInvalidRouterParamError(paramName string) InvalidRouterParamError {
	return InvalidRouterParamError{
		ParamName: paramName,
	}
}

// Error returns the error message.
func (e InvalidRouterParamError) Error() string {
	return "invalid router param: " + e.ParamName
}
