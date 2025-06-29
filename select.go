package database

import (
	"errors"

	"github.com/gouniverse/maputils"
)

// SelectToMapAny executes a SQL query in the given context and returns a slice of maps,
// where each map represents a row of the query results. The keys of the map are the
// column names of the query, and the values are the values of the columns.
//
// The context is used to control the execution of the query, allowing for
// cancellation and timeout control. It also allows to be used with DB, Tx, and Conn.
//
// If the query returns no rows, the function returns an empty slice.
//
// Example usage:
//
// listMap, err := SelectToMapAny(context.Background(), "SELECT * FROM users")
//
// Parameters:
// - ctx (context.Context): The context to use for the query execution.
// - sqlStr (string): The SQL query to execute.
// - args (any): Optional arguments to pass to the query.
//
// Returns:
// - []map[string]any: A slice of maps containing the query results.
// - error: An error if the query failed.
func SelectToMapAny(ctx QueryableContext, sqlStr string, args ...any) ([]map[string]any, error) {
	if ctx.queryable == nil {
		return []map[string]any{}, errors.New("querier (db/tx/conn) is nil")
	}

	listMap := []map[string]any{}

	rows, err := ctx.queryable.QueryContext(ctx, sqlStr, args...)

	if err != nil {
		return []map[string]any{}, err
	}
	defer rows.Close()

	// Get column names
	columns, err := rows.Columns()
	if err != nil {
		return []map[string]any{}, err
	}

	for rows.Next() {
		// Create a slice of interface{} to hold the values
		values := make([]interface{}, len(columns))
		// Create a slice of pointers to interface{} for scanning
		valuePtrs := make([]interface{}, len(columns))
		for i := range values {
			valuePtrs[i] = &values[i]
		}

		// Scan the row into the slice of pointers
		if err := rows.Scan(valuePtrs...); err != nil {
			return []map[string]any{}, err
		}

		// Create a map for this row
		row := make(map[string]any)
		for i, col := range columns {
			val := values[i]
			// Handle nil values
			if val == nil {
				row[col] = nil
			} else {
				row[col] = val
			}
		}

		listMap = append(listMap, row)
	}

	if err := rows.Err(); err != nil {
		return []map[string]any{}, err
	}

	return listMap, nil
}

// SelectToMapString executes a SQL query in the given context and returns a slice of maps,
// where each map represents a row of the query results. The keys of the map are the
// column names of the query, and the values are the values of the columns as strings.
//
// The context is used to control the execution of the query, allowing for
// cancellation and timeout control. It also allows to be used with DB, Tx, and Conn.
//
// If the query returns no rows, the function returns an empty slice.
//
// Example usage:
//
// listMap, err := SelectToMapString(context.Background(), "SELECT * FROM users")
//
// Parameters:
// - ctx (context.Context): The context to use for the query execution.
// - sqlStr (string): The SQL query to execute.
// - args (any): Optional arguments to pass to the query.
//
// Returns:
// - []map[string]string: A slice of maps containing the query results.
// - error: An error if the query failed.
func SelectToMapString(ctx QueryableContext, sqlStr string, args ...any) ([]map[string]string, error) {
	if ctx.queryable == nil {
		return []map[string]string{}, errors.New("querier (db/tx/conn) is nil")
	}

	listMapAny, err := SelectToMapAny(ctx, sqlStr, args...)

	if err != nil {
		return []map[string]string{}, err
	}

	listMapString := []map[string]string{}

	// Iterate over the list of maps and convert each map from map[string]any to map[string]string.
	// This is done by using the maputils.MapStringAnyToMapStringString() function.
	for i := range listMapAny {
		mapString := maputils.MapStringAnyToMapStringString(listMapAny[i])
		listMapString = append(listMapString, mapString)
	}

	return listMapString, nil
}
