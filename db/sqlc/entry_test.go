package db

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"gitlab.com/AdamClaassens/simple_bank/util"
)

// Crate random entry and pass in the account
func createRandomEntry(t *testing.T, account Account) Entry {
	// Create arguments to use for the entry
	arg := CreateEntryParams{
		AccountID: account.ID,
		Amount:    util.RandomMoney(),
	}

	// Create the entry
	entry, err := testQueries.CreateEntry(context.Background(), arg)
	// Check for errors
	require.NoError(t, err)
	// Check if the entry is not empty
	require.NotEmpty(t, entry)

	// Check if the arguments used for account_id and amount are the same as the returned entry
	require.Equal(t, arg.AccountID, entry.AccountID)
	// Check if the arguments used for amount and amount are the same as the returned entry
	require.Equal(t, arg.Amount, entry.Amount)

	// Check if the returned entry id and created_at are not zero
	require.NotZero(t, entry.ID)
	// Check if the returned entry created_at is not zero
	require.NotZero(t, entry.CreatedAt)

	// Return the created entry
	return entry
}

// TestCreateEntry tests the CreateEntry function
func TestCreateEntry(t *testing.T) {
	// Create a random account
	account := createRandomAccount(t)
	// Create a random entry and pass in the account
	createRandomEntry(t, account)
}

// TestGetEntry tests the GetEntry function
func TestGetEntry(t *testing.T) {
	// Create a random account
	account := createRandomAccount(t)
	// Create a random entry and pass in the account
	entry1 := createRandomEntry(t, account)
	// Get the entry using the entry id
	entry2, err := testQueries.GetEntry(context.Background(), entry1.ID)
	// Check for errors
	require.NoError(t, err)
	// Check if the entry is not empty
	require.NotEmpty(t, entry2)

	// Check if the created entry and the entry from get have the same ID
	require.Equal(t, entry1.ID, entry2.ID)
	// Check if the created entry and the entry from get have the same AccountID
	require.Equal(t, entry1.AccountID, entry2.AccountID)
	// Check if the created entry and the entry from get have the same Amount
	require.Equal(t, entry1.Amount, entry2.Amount)
	// Check if the created entry and the entry from get have the same CreatedAt
	require.WithinDuration(t, entry1.CreatedAt, entry2.CreatedAt, time.Second)
}

// TestListEntries tests the ListEntries function
func TestListEntries(t *testing.T) {
	// Create a random account
	account := createRandomAccount(t)
	// Create 10 random entries and pass in the account
	for i := 0; i < 10; i++ {
		createRandomEntry(t, account)
	}

	// Create arguments to use for the ListEntries function
	arg := ListEntriesParams{
		AccountID: account.ID,
		Limit:     5,
		Offset:    5,
	}

	// List the entries using the above arguments
	entries, err := testQueries.ListEntries(context.Background(), arg)
	// Check for errors
	require.NoError(t, err)
	// Check that 5 entries were returned
	require.Len(t, entries, 5)

	// Loop through the returned entries
	for _, entry := range entries {
		// Check if the entry is not empty
		require.NotEmpty(t, entry)
		//Check if the arguments used for account_id and amount are the same as the returned entry
		require.Equal(t, arg.AccountID, entry.AccountID)
	}
}
