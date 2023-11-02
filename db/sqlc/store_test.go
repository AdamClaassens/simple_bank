package db

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/require"
	"math"
	"testing"
)

func TestStore_TransferTx(t *testing.T) {
	// Create a new store Object with the test database
	store := NewStore(testDB)

	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)
	fmt.Println(">> Before:", account1.Balance, account2.Balance)

	// Create a random transfer and pass in the accounts
	n := 5
	amount := 10.0

	errs := make(chan error)
	results := make(chan TransferTxResult)

	for i := 0; i < n; i++ {
		//txName := fmt.Sprintf("tx %d", i+1)
		go func() {
			//ctx := context.WithValue(context.Background(), txKey, txName)
			//result, err := store.TransferTX(ctx, TransferTxParams{
			result, err := store.TransferTX(context.Background(), TransferTxParams{
				FromAccountID: account1.ID,
				ToAccountID:   account2.ID,
				Amount:        amount,
			})
			errs <- err
			results <- result
		}()
	}

	// Check if there are no errors
	existed := make(map[int]bool)
	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)

		result := <-results
		require.NotEmpty(t, result)

		transfer := result.Transfer
		// Check if the transfer is not empty
		require.NotEmpty(t, result.Transfer)
		// Check if the transfer from_account_id is the same as the account1 id
		require.Equal(t, account1.ID, result.Transfer.FromAccountID)
		// Check if the transfer to_account_id is the same as the account2 id
		require.Equal(t, account2.ID, result.Transfer.ToAccountID)
		// Check if the transfer amount is the same as the amount
		require.Equal(t, amount, result.Transfer.Amount)
		// Check if the transfer id is not zero
		require.NotZero(t, result.Transfer.ID)
		// Check if the transfer created_at is not zero
		require.NotZero(t, result.Transfer.CreatedAt)

		_, err = store.GetTransfer(context.Background(), transfer.ID)
		require.NoError(t, err)

		fromEntry := result.FromEntry
		// Check if the fromEntry is not empty
		require.NotEmpty(t, fromEntry)
		// Check if the fromEntry account_id is the same as the account1 id
		require.Equal(t, account1.ID, fromEntry.AccountID)
		// Check if the fromEntry amount is the same as the amount
		require.Equal(t, -amount, fromEntry.Amount)
		// Check if the fromEntry id is not zero
		require.NotZero(t, fromEntry.ID)
		// Check if the fromEntry created_at is not zero
		require.NotZero(t, fromEntry.CreatedAt)

		_, err = store.GetEntry(context.Background(), fromEntry.ID)
		require.NoError(t, err)

		toEntry := result.ToEntry
		// Check if the toEntry is not empty
		require.NotEmpty(t, toEntry)
		// Check if the toEntry account_id is the same as the account2 id
		require.Equal(t, account2.ID, toEntry.AccountID)
		// Check if the toEntry amount is the same as the amount
		require.Equal(t, amount, toEntry.Amount)
		// Check if the toEntry id is not zero
		require.NotZero(t, toEntry.ID)
		// Check if the toEntry created_at is not zero
		require.NotZero(t, toEntry.CreatedAt)

		_, err = store.GetEntry(context.Background(), toEntry.ID)
		require.NoError(t, err)

		// Check accounts
		fromAccount := result.FromAccount
		// Check if the fromAccount is not empty
		require.NotEmpty(t, fromAccount)
		// Check if the fromAccount id is the same as the account1 id
		require.Equal(t, account1.ID, fromAccount.ID)

		toAccount := result.ToAccount
		// Check if the toAccount is not empty
		require.NotEmpty(t, toAccount)
		// Check if the toAccount id is the same as the account2 id
		require.Equal(t, account2.ID, toAccount.ID)

		fmt.Println(">> From:", fromAccount.Balance, "To:", toAccount.Balance)
		// Check account balances
		diff1 := roundFloat(account1.Balance-fromAccount.Balance, 2)
		diff2 := roundFloat(toAccount.Balance-account2.Balance, 2)
		// Check if the account1 balance is the same as the amount
		require.Equal(t, diff1, diff2)
		require.True(t, diff1 > 0)
		require.True(t, math.Mod(diff1, amount) == 0) // 1 * amount, 2 * amount, 3 * amount, ..., n * amount

		k := int(diff1 / amount)
		require.True(t, k >= 1 && int(k) <= n)
		require.NotContains(t, existed, k)
		existed[k] = true
	}

	// Check the account balances after the transfer
	updateAccount1, err := testQueries.GetAccount(context.Background(), account1.ID)
	// Check if there are no errors
	require.NoError(t, err)

	// Check the account balances after the transfer
	updateAccount2, err := testQueries.GetAccount(context.Background(), account2.ID)
	// Check if there are no errors
	require.NoError(t, err)

	fmt.Println(">> After:", updateAccount1.Balance, updateAccount2.Balance)
	// Check if the account1 balance is the same as the amount
	require.Equal(t, account1.Balance-float64(n)*amount, updateAccount1.Balance)
	// Check if the account2 balance is the same as the amount
	require.Equal(t, account2.Balance+float64(n)*amount, updateAccount2.Balance)
}

func TestStore_TransferTxDeadlock(t *testing.T) {
	// Create a new store Object with the test database
	store := NewStore(testDB)

	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)
	fmt.Println(">> Before:", account1.Balance, account2.Balance)

	// Create a random transfer and pass in the accounts
	n := 10
	amount := 10.0

	errs := make(chan error)

	for i := 0; i < n; i++ {
		fromAccountID := account1.ID
		toAccountID := account2.ID
		if i%2 == 1 {
			fromAccountID = account2.ID
			toAccountID = account1.ID
		}

		//txName := fmt.Sprintf("tx %d", i+1)
		go func() {
			//ctx := context.WithValue(context.Background(), txKey, txName)
			//result, err := store.TransferTX(ctx, TransferTxParams{
			_, err := store.TransferTX(context.Background(), TransferTxParams{
				FromAccountID: fromAccountID,
				ToAccountID:   toAccountID,
				Amount:        amount,
			})
			errs <- err
		}()
	}
	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)
	}

	// Check the account balances after the transfer
	updateAccount1, err := testQueries.GetAccount(context.Background(), account1.ID)
	// Check if there are no errors
	require.NoError(t, err)

	// Check the account balances after the transfer
	updateAccount2, err := testQueries.GetAccount(context.Background(), account2.ID)
	// Check if there are no errors
	require.NoError(t, err)

	fmt.Println(">> After:", updateAccount1.Balance, updateAccount2.Balance)
	// Chekc if account 1 balance is the same as it was before since we are transferring the same amount back and forth
	require.Equal(t, account1.Balance, updateAccount1.Balance)
	// check if account 2 balance is the same as it was before since we are transferring the same amount back and forth
	require.Equal(t, account2.Balance, updateAccount2.Balance)
}
