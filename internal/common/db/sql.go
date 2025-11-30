package db

import (
	"database/sql"
	"fmt"
)

// CheckTransaction handles the commit or rollback of a database transaction
// based on the state of the provided error pointer.
//
// If the error pointer (*repoError) is not nil, the function attempts to roll back
// the transaction. If the rollback fails, it wraps the rollback error with the
// existing error.
//
// If the error pointer is nil, the function attempts to commit the transaction.
// If the commit fails, it sets the error pointer to a wrapped error containing
// the commit failure details.
//
// Parameters:
//   - tx: A pointer to the sql.Tx object representing the database transaction.
//   - repoError: A pointer to an error that determines whether to commit or roll back
//     the transaction and is updated with any errors encountered during the process.
func CheckTransaction(tx *sql.Tx, repoError *error) {
	if *repoError != nil {
		if err := tx.Rollback(); err != nil {
			*repoError = fmt.Errorf("(repo) failed to rollback transaction: %w: %w", err, *repoError)
		}
	} else {
		if err := tx.Commit(); err != nil {
			*repoError = fmt.Errorf("(repo) failed to commit transaction: %w: %w", err, *repoError)
		}
	}
}
