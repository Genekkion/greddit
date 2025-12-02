package postgres

import (
	"greddit/internal/test"
	"testing"
)

func TestCreateTransaction(t *testing.T) {
	t.Parallel()

	pool, cleanup := NewTestPool(t)
	defer cleanup()

	txs := NewTransactional(pool)

	ctx := t.Context()

	var err error
	ctx, err = txs.CtxTx(ctx)
	test.NilErr(t, err)

	tx, err := txs.ctxGetTx(ctx)
	test.NilErr(t, err)

	test.Assert(t, "Transaction should have been created", tx != nil)

	err = tx.Commit(ctx)
	test.NilErr(t, err)
}

func TestGetTxFromNilContext(t *testing.T) {
	t.Parallel()

	pool, cleanup := NewTestPool(t)
	defer cleanup()

	txs := NewTransactional(pool)

	ctx, err := txs.CtxTx(nil)
	test.Assert(t, "Expect error from being nil context", err != nil)
	test.Assert(t, "Expect nil context", ctx == nil)
}
