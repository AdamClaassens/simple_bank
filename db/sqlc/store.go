package db

import (
	"context"
	"database/sql"
	"fmt"
	"math"
)

// Round a float value to a certain precision
func roundFloat(val float64, precision uint) float64 {
	ratio := math.Pow(10, float64(precision))
	return math.Round(val*ratio) / ratio
}

// Store provides all functions to execute db queries and transactions
type Store struct {
	*Queries
	db *sql.DB
}

// NewStore creates a new store
func NewStore(db *sql.DB) *Store {
	return &Store{
		db:      db,
		Queries: New(db),
	}
}

// execTx executes a function within a database transaction
func (store *Store) execTx(ctx context.Context, fn func(*Queries) error) error {
	// Create a new transaction
	tx, err := store.db.BeginTx(ctx, nil)
	// Check for errors
	if err != nil {
		return err
	}

	// Create a new queries struct with the transaction
	q := New(tx)
	// Execute the function
	err = fn(q)
	// Check for errors
	if err != nil {
		// Rollback the transaction
		if rbErr := tx.Rollback(); rbErr != nil {
			// Return the error
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}
		// Return the error
		return err
	}

	// Commit the transaction
	return tx.Commit()
}

// TransferTxParams contains the input parameters for the Transfer transaction
type TransferTxParams struct {
	FromAccountID int64   `json:"from_account_id"`
	ToAccountID   int64   `json:"to_account_id"`
	Amount        float64 `json:"amount"`
}

// TransferTxResult is the result of the TransferTx function
type TransferTxResult struct {
	// The newly created transfer
	Transfer Transfer `json:"transfer"`
	// The entry recorded for the account that the transfer was made from
	FromAccount Account `json:"from_account"`
	// The entry recorded for the account that the transfer was made to
	ToAccount Account `json:"to_account"`
	// The entry recorded for the account that the transfer was made from
	FromEntry Entry `json:"from_entry"`
	// The entry recorded for the account that the transfer was made to
	ToEntry Entry `json:"to_entry"`
}

// TransferTx performs a money transfer from one account to the other
// It creates a transfer record, add account entries, and update accounts' balance within a single database transaction
func (store *Store) TransferTX(ctx context.Context, arg TransferTxParams) (TransferTxResult, error) {
	var result TransferTxResult

	// Execute the transaction
	err := store.execTx(ctx, func(q *Queries) error {
		// Create a new transfer
		var err error

		// Create a new transfer
		result.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams{
			FromAccountID: arg.FromAccountID,
			ToAccountID:   arg.ToAccountID,
			Amount:        arg.Amount,
		})
		// Check for errors
		if err != nil {
			return err
		}

		// Create a new entry for the account that the money is moving out of
		result.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.FromAccountID,
			Amount:    -arg.Amount,
		})
		// Check for errors
		if err != nil {
			return err
		}

		// Create a new entry for the account that the money is moving into
		result.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.ToAccountID,
			Amount:    arg.Amount,
		})
		// Check for errors
		if err != nil {
			return err
		}

		// Check if the from account id is less than the to account id
		// Then update the from account balance first
		// This is to prevent deadlocks
		if arg.FromAccountID < arg.ToAccountID {
			// Update the from account balance first
			result.FromAccount, result.ToAccount, err = addMoney(ctx, q, arg.FromAccountID, -arg.Amount, arg.ToAccountID, arg.Amount)
		} else {
			// Update the to account balance first
			result.ToAccount, result.FromAccount, err = addMoney(ctx, q, arg.ToAccountID, arg.Amount, arg.FromAccountID, -arg.Amount)

		}

		// Return nil for the error
		return nil
	})

	// Return the result and error
	return result, err
}

// addMoney updates the account balance by a certain amount and returns the updated account
func addMoney(
	ctx context.Context,
	q *Queries,
	accountID1 int64,
	amount1 float64,
	accountID2 int64,
	amount2 float64,
) (account1 Account, account2 Account, err error) {
	// Get the account after the balance has been updated
	account1, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
		ID:     accountID1,
		Amount: amount1,
	})
	if err != nil {
		// Return the error
		return
	}

	// Get the account after the balance has been updated
	account2, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
		ID:     accountID2,
		Amount: amount2,
	})
	if err != nil {
		// Return the error
		return
	}

	// Return the accounts
	return
}
