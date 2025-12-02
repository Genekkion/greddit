package dbports

import "context"

// Transactional is for database transactions, enabling nearly invisible usage of
// transactions by piggybacking on a transaction within the context.
type Transactional interface {
	// CtxTx returns a context with a transaction. This transaction should be
	// invisible to the users of this interface and should be extracted from the
	// context by the underlying implementation.
	CtxTx(ctx context.Context) (context.Context, error)

	// TxRollback rolls back the transaction.
	TxRollback(ctx context.Context) error

	// TxCommit commits the transaction.
	TxCommit(ctx context.Context) error
}
