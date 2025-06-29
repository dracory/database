package database

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	_ "modernc.org/sqlite"
)

// initSqliteDB creates a new SQLite database connection in memory.
//
// Returns:
// - *sql.DB: the database connection
// - error: the error if any
func initSqliteDB() (*sql.DB, error) {
	// Create a new database connection in memory
	db, err := Open(Options().
		SetDatabaseType(DATABASE_TYPE_SQLITE).
		SetDatabaseHost("").
		SetDatabasePort("").
		SetDatabaseName(":memory:").
		SetUserName("").
		SetPassword(""))

	// Check if there was an error
	if err != nil {
		return nil, err
	}

	// Check if the database connection is not nil
	if db == nil {
		return nil, errors.New("db is nil")
	}

	// Return the database connection
	return db, nil
}

func TestDatabaseTypeFromDB(t *testing.T) {
	db, err := initSqliteDB()

	if err != nil {
		t.Fatal(err)
	}

	dbType := DatabaseType(db)

	if dbType != DATABASE_TYPE_SQLITE {
		t.Fatalf("Expected Debug [%v], received [%v]", "sqlite", dbType)
	}
}

func TestDatabaseTypeFromTx(t *testing.T) {
	db, err := initSqliteDB()

	if err != nil {
		t.Fatal(err)
	}

	tx, err := db.Begin()

	if err != nil {
		t.Fatal(err)
	}

	dbType := DatabaseType(tx)

	if dbType != DATABASE_TYPE_SQLITE {
		t.Fatalf("Expected Debug [%v], received [%v]", "sqlite", dbType)
	}
}

func TestDatabaseTypeFromConn(t *testing.T) {
	db, err := initSqliteDB()

	if err != nil {
		t.Fatal(err)
	}

	conn, err := db.Conn(context.Background())

	if err != nil {
		t.Fatal(err)
	}

	dbType := DatabaseType(conn)

	if dbType != DATABASE_TYPE_SQLITE {
		t.Fatalf("Expected Debug [%v], received [%v]", "sqlite", dbType)
	}
}
