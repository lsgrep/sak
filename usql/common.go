package usql

import (
	"context"
	"database/sql"
	"errors"
	"github.com/lsgrep/sak/utime"
)

var (
	NotFound = errors.New("Not Found")
)

// Implemented by both `*sql.Row` and `*sql.Rows`.
type RowScanner interface {
	Scan(dest ...interface{}) error
}

// Implemented by both `*sql.DB`, `*sql.Stmt` and `*sql.Tx`.
type Querier interface {
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
}

func newQueryCtx() context.Context {
	return utime.CtxTimeoutMs(5e3)
}
