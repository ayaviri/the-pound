package db

import "database/sql"

var err error

type DBExecutor interface {
	Exec(query string, args ...any) (sql.Result, error)
	Query(query string, args ...any) (*sql.Rows, error)
	QueryRow(query string, args ...any) *sql.Row
}

// Expects a raw SQL count query, writing the first column of the first row
// returned into a string, converting it into an integer, and returning it.
// Returns any errors encountered along the way as well
func QueryCount(e DBExecutor, query string, args ...any) (int, error) {
	var count int
	var row *sql.Row = e.QueryRow(query, args...)
	err = row.Scan(&count)

	if err != nil {
		return 0, err
	}

	return count, nil
}

func ExecuteInTransaction(db *sql.DB, operation func(DBExecutor) error) error {
	tx, err := db.Begin()

	if err != nil {
		return err
	}

	defer tx.Rollback()

	if err = operation(tx); err != nil {
		return err
	}

	return tx.Commit()
}
