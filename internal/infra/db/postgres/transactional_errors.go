package postgres

var (
	NoTxInCtxError      = noTxInCtxError{}
	InvalidTxInCtxError = invalidTxInCtxError{}
	NilCtxError         = nilCtxError{}
)

// noTxInCtxError represents an error when there is no transaction in the context.
type noTxInCtxError struct{}

// Error returns the error message.
func (e noTxInCtxError) Error() string {
	return "no transaction in context"
}

// invalidTxInCtxError represents an error when the transaction in the context is invalid.
type invalidTxInCtxError struct{}

// Error returns the error message.
func (e invalidTxInCtxError) Error() string {
	return "invalid transaction in context"
}

// nilCtxError represents an error when the context is nil.
type nilCtxError struct{}

// Error returns the error message.
func (e nilCtxError) Error() string {
	return "nil context"
}
