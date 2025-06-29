package database

import (
	"database/sql"
	"errors"
)

// Query executes a SQL query in the given context and returns a *sql.Rows object containing the query results.
//
// The context is used to control the execution of the query, allowing
// for cancellation and timeout control. It also allows to be used with
// DB, Tx, and Conn.
//
// Example usage:
//
// rows, err := Query(context.Background(), "SELECT * FROM users")
//
// Parameters:
// - ctx (context.Context): The context to use for the query execution.
// - sqlStr (string): The SQL query to execute.
// - args (any): Optional arguments to pass to the query.
//
// Returns:
// - *sql.Rows: A *sql.Rows object containing the query results.
// - error: An error if the query failed.
func Query(ctx QueryableContext, sqlStr string, args ...any) (*sql.Rows, error) {
	// Check for nil querier
	if ctx.queryable == nil {
		return nil, errors.New("querier (db/tx/conn) is nil")
	}

	// Ensure the context is properly wrapped with the queryable
	ctx = NewQueryableContextOr(ctx, ctx.queryable)

	// Execute the query in the context
	return ctx.queryable.QueryContext(ctx, sqlStr, args...)
}
