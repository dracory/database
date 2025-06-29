package database

import (
	"context"
	"database/sql"
	"testing"
)

// TestTransactionWithContextOr tests the use of transactions with ContextOr
func TestTransactionWithContextOr(t *testing.T) {
	// Open an in-memory SQLite database for testing
	// Note: For SQLite in-memory databases, we need to use a shared cache to test transactions
	// properly across different connections
	db, err := sql.Open("sqlite", "file::memory:?cache=shared")
	if err != nil {
		t.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	// Create a test table
	_, err = db.Exec(`CREATE TABLE users (
		id INTEGER PRIMARY KEY,
		name TEXT NOT NULL,
		email TEXT NOT NULL
	)`)
	if err != nil {
		t.Fatalf("Failed to create table: %v", err)
	}

	// Begin a transaction
	tx, err := db.BeginTx(context.Background(), nil)
	if err != nil {
		t.Fatalf("Failed to begin transaction: %v", err)
	}

	// Create a transaction context using Context
	txCtx := Context(context.Background(), tx)

	// Test function that simulates a data store method using ContextOr
	simulateStoreMethod := func(ctx context.Context, db *sql.DB, name, email string) error {
		// Convert to queryable context, preserving any existing transaction
		qCtx := ContextOr(ctx, db)

		// Execute the query using the queryable context
		_, err := Execute(qCtx, "INSERT INTO users (name, email) VALUES (?, ?)", name, email)
		return err
	}

	// Test 1: Call the function with a transaction context
	err = simulateStoreMethod(txCtx, db, "John Doe", "john@example.com")
	if err != nil {
		tx.Rollback()
		t.Fatalf("Failed to insert data with transaction context: %v", err)
	}

	// For Test 2, we'll use the transaction context again
	// This demonstrates that multiple operations can be part of the same transaction
	err = simulateStoreMethod(txCtx, db, "Jane Doe", "jane@example.com")
	if err != nil {
		tx.Rollback()
		t.Fatalf("Failed to insert data with transaction context: %v", err)
	}

	// Verify that the inserts are in the transaction by querying within the transaction
	var countInTransaction int
	err = tx.QueryRow("SELECT COUNT(*) FROM users").Scan(&countInTransaction)
	if err != nil {
		tx.Rollback()
		t.Fatalf("Failed to count rows in transaction: %v", err)
	}

	// We should see both rows in the transaction
	if countInTransaction != 2 {
		tx.Rollback()
		t.Errorf("Expected 2 rows in transaction, got %d", countInTransaction)
	}

	// Commit the transaction
	err = tx.Commit()
	if err != nil {
		t.Fatalf("Failed to commit transaction: %v", err)
	}

	// Now check the row count after commit
	var countAfterCommit int
	err = db.QueryRow("SELECT COUNT(*) FROM users").Scan(&countAfterCommit)
	if err != nil {
		t.Fatalf("Failed to count rows after commit: %v", err)
	}

	// We should see both rows now
	if countAfterCommit != 2 {
		t.Errorf("Expected 2 rows after commit, got %d", countAfterCommit)
	}
}

// TestTransactionRollback tests that changes are properly rolled back when a transaction is aborted
func TestTransactionRollback(t *testing.T) {
	// Open an in-memory SQLite database for testing
	// Using shared cache for consistent behavior across connections
	db, err := sql.Open("sqlite", "file::memory:?cache=shared")
	if err != nil {
		t.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	// Create a test table
	_, err = db.Exec(`CREATE TABLE products (
		id INTEGER PRIMARY KEY,
		name TEXT NOT NULL,
		price REAL NOT NULL
	)`)
	if err != nil {
		t.Fatalf("Failed to create table: %v", err)
	}

	// Begin a transaction
	tx, err := db.BeginTx(context.Background(), nil)
	if err != nil {
		t.Fatalf("Failed to begin transaction: %v", err)
	}

	// Create a transaction context
	txCtx := Context(context.Background(), tx)

	// Insert some data in the transaction
	_, err = Execute(txCtx, "INSERT INTO products (name, price) VALUES (?, ?)", "Product 1", 19.99)
	if err != nil {
		tx.Rollback()
		t.Fatalf("Failed to insert data: %v", err)
	}

	// Simulate an error condition
	// In a real application, we'd have validation that would cause an error for negative prices
	// For this test, we'll just simulate the error and roll back
	err = tx.Rollback()
	if err != nil {
		t.Fatalf("Failed to rollback transaction: %v", err)
	}

	// Create a new transaction for verification
	// This is necessary because the previous transaction is now closed
	verifyTx, err := db.BeginTx(context.Background(), nil)
	if err != nil {
		t.Fatalf("Failed to begin verification transaction: %v", err)
	}
	defer verifyTx.Rollback()

	// Check that no rows were committed
	var count int
	err = verifyTx.QueryRow("SELECT COUNT(*) FROM products").Scan(&count)
	if err != nil {
		t.Fatalf("Failed to count rows: %v", err)
	}

	// We should see no rows since the transaction was rolled back
	if count != 0 {
		t.Errorf("Expected 0 rows after rollback, got %d", count)
	}
}
