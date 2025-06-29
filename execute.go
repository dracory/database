package database

import (
	"database/sql"
	"errors"
)

// Execute executes a SQL query in the given context and returns a sql.Result
// containing information about the execution, or an error if the query failed.
//
// The context is used to control the execution of the query, allowing for
// cancellation and timeout control. It also allows to be used with
// DB, Tx, and Conn.
//
// Example usage:
//
// result, err := Execute(context.Background(), "UPDATE users SET name = ? WHERE id = ?", "John Doe", 1)
//
// Parameters:
// - ctx (context.Context): The context to use for the query execution.
// - sqlStr (string): The SQL query to execute.
// - args (any): Optional arguments to pass to the query.
//
// Returns:
// - sql.Result: A sql.Result object containing information about the execution.
// - error: An error if the query failed.
func Execute(ctx QueryableContext, sqlStr string, args ...any) (sql.Result, error) {
	// Check if the querier is nil
	if ctx.queryable == nil {
		return nil, errors.New("querier (db/tx/conn) is nil")
	}

	ctx = NewQueryableContextOr(ctx, ctx.queryable)

	// Execute the query
	return ctx.queryable.ExecContext(ctx, sqlStr, args...)
}
