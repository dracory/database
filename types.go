package database

import (
	"context"
	"database/sql"
)

// NewQueryableContext returns a new context with the given QueryableInterface.
//
// Note: For convenience, a shortcut alias function 'Context' is provided in funcs.go
// that calls this function with the same parameters.
func NewQueryableContext(ctx context.Context, queryable QueryableInterface) QueryableContext {
	return QueryableContext{Context: ctx, queryable: queryable}
}

// NewQueryableContextOr returns the existing QueryableContext if the provided context
// is already a QueryableContext, or creates a new one with the given QueryableInterface.
//
// Note: For convenience, a shortcut alias function 'ContextOr' is provided in funcs.go
// that calls this function with the same parameters.
func NewQueryableContextOr(ctx context.Context, queryable QueryableInterface) QueryableContext {
	if qCtx, ok := ctx.(QueryableContext); ok {
		return qCtx
	}

	return QueryableContext{Context: ctx, queryable: queryable}
}

// Verify that QueryableContext implements the context.Context interface.
var _ context.Context = QueryableContext{}

// QueryableContext extends the context.Context interface with a queryable field.
// The queryable field may be of type *sql.DB, *sql.Conn, or *sql.Tx.
type QueryableContext struct {
	context.Context
	queryable QueryableInterface
}

func (ctx QueryableContext) IsDB() bool {
	if ctx.queryable == nil {
		return false
	}

	if ctx.IsTx() {
		return false
	}

	if ctx.IsConn() {
		return false
	}

	_, ok := ctx.queryable.(*sql.DB)

	return ok
}

func (ctx QueryableContext) IsConn() bool {
	if ctx.queryable == nil {
		return false
	}

	_, ok := ctx.queryable.(*sql.Conn)

	return ok
}

func (ctx QueryableContext) IsTx() bool {
	if ctx.queryable == nil {
		return false
	}

	_, ok := ctx.queryable.(*sql.Tx)

	return ok
}

func (ctx QueryableContext) Queryable() QueryableInterface {
	return ctx.queryable
}
