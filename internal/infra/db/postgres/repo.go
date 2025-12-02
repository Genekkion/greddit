package postgres

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

// BaseRepo is the base repository for all repositories, based on pgxpool.Pool.
type BaseRepo struct {
	db  *pgxpool.Pool
	txs Transactional
}

// NewBaseRepo creates a new BaseRepo.
func NewBaseRepo(pool *pgxpool.Pool) BaseRepo {
	return BaseRepo{
		db: pool,
		txs: Transactional{
			db: pool,
		},
	}
}

// ErrRow is a pgx.Row that returns an error on Scan.
type ErrRow struct {
	err error
}

// Scan returns the error.
func (e ErrRow) Scan(_ ...any) error {
	return e.err
}

// QueryRow executes a query expected to return at most one row. It uses the
// transaction in the context if available, otherwise it uses the pool directly.
func (r *BaseRepo) QueryRow(ctx context.Context, stmt string, args ...any) (row pgx.Row) {
	tx, err := r.txs.ctxGetTx(ctx)
	if err != nil {
		if errors.Is(err, NoTxInCtxError) {
			return r.db.QueryRow(ctx, stmt, args...)
		} else {
			return ErrRow{
				err: err,
			}
		}
	}

	return tx.QueryRow(ctx, stmt, args...)
}

// Query executes a query that returns rows. It uses the transaction in the
// context if available, otherwise it uses the pool directly.
func (r *BaseRepo) Query(ctx context.Context, stmt string, args ...any) (rows pgx.Rows, err error) {
	tx, err := r.txs.ctxGetTx(ctx)
	if err != nil {
		if errors.Is(err, NoTxInCtxError) {
			return r.db.Query(ctx, stmt, args...)
		} else {
			return nil, err
		}
	}

	return tx.Query(ctx, stmt, args...)
}

// Exec executes an SQL statement. It uses the transaction in the context if
// available, otherwise it uses the pool directly.
func (r *BaseRepo) Exec(ctx context.Context, stmt string, args ...any) (commandTag pgconn.CommandTag, err error) {
	tx, err := r.txs.ctxGetTx(ctx)
	if err != nil {
		if errors.Is(err, NoTxInCtxError) {
			return r.db.Exec(ctx, stmt, args...)
		} else {
			return pgconn.CommandTag{}, err
		}
	}
	return tx.Exec(ctx, stmt, args...)
}
