package db

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"gitlab.com/AdamClaassens/simple_bank/util"
)

// Crate random transfer and pass in the accounts
func createRandomTransfer(t *testing.T, account1, account2 Account) Transfer {
	// Create arguments to use for the transfer
	arg := CreateTransferParams{
		FromAccountID: account1.ID,
		ToAccountID:   account2.ID,
		Amount:        util.RandomMoney(),
	}

	// Create the transfer
	transfer, err := testQueries.CreateTransfer(context.Background(), arg)
	// Check for errors
	require.NoError(t, err)
	// Check if the transfer is not empty
	require.NotEmpty(t, transfer)

	// Check if the arguments used to create transfer from_account_id and from_account_id are the same as the returned transfer
	require.Equal(t, arg.FromAccountID, transfer.FromAccountID)
	// Check if the arguments used to create transfer to_account_id and to_account_id are the same as the returned transfer
	require.Equal(t, arg.ToAccountID, transfer.ToAccountID)
	// Check if the arguments used to create transfer amount and amount are the same as the returned transfer
	require.Equal(t, arg.Amount, transfer.Amount)

	// Check if the returned transfer id is not zero
	require.NotZero(t, transfer.ID)
	// Check if the returned transfer created_at is not zero
	require.NotZero(t, transfer.CreatedAt)

	// Return the created transfer
	return transfer
}

// TestCreateTransfer tests the CreateTransfer function
func TestCreateTransfer(t *testing.T) {
	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)
	// Create a random transfer and pass in the accounts
	createRandomTransfer(t, account1, account2)
}

// TestGetTransfer tests the GetTransfer function
func TestGetTransfer(t *testing.T) {
	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)
	transfer1 := createRandomTransfer(t, account1, account2)

	// Get the transfer using the transfer id
	transfer2, err := testQueries.GetTransfer(context.Background(), transfer1.ID)
	// Check for errors
	require.NoError(t, err)
	// Check if the transfer is not empty
	require.NotEmpty(t, transfer2)

	// Check if the returned transfer id is the same as the created transfer id
	require.Equal(t, transfer1.ID, transfer2.ID)
	// Check if the returned transfer from_account_id is the same as the created transfer from_account_id
	require.Equal(t, transfer1.FromAccountID, transfer2.FromAccountID)
	// Check if the returned transfer to_account_id is the same as the created transfer to_account_id
	require.Equal(t, transfer1.ToAccountID, transfer2.ToAccountID)
	// Check if the returned transfer amount is the same as the created transfer amount
	require.Equal(t, transfer1.Amount, transfer2.Amount)
	// Check if the returned transfer created_at is the same as the created transfer created_at
	require.WithinDuration(t, transfer1.CreatedAt, transfer2.CreatedAt, time.Second)
}

// TestListTransfer tests the ListTransfer function
func TestListTransfer(t *testing.T) {
	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)

	// Create 10 random transfers and pass in the accounts
	for i := 0; i < 5; i++ {
		createRandomTransfer(t, account1, account2)
		createRandomTransfer(t, account2, account1)
	}

	// Create arguments to use for the list transfers
	arg := ListTransfersParams{
		FromAccountID: account1.ID,
		ToAccountID:   account1.ID,
		Limit:         5,
		Offset:        5,
	}

	// List the transfers
	transfers, err := testQueries.ListTransfers(context.Background(), arg)
	// Check for errors
	require.NoError(t, err)
	// Check that 5 transfers were returned
	require.Len(t, transfers, 5)

	// Loop through the returned transfers
	for _, transfer := range transfers {
		// Check if the transfer is not empty
		require.NotEmpty(t, transfer)
		// Check if the arguments used for from_account_id and to_account_id are the same as the returned transfer
		require.True(t, transfer.FromAccountID == account1.ID || transfer.ToAccountID == account1.ID)
	}
}
