package db

import (
	"context"
	"simplebank/util"
	"testing"
	"time"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/require"
)

func TestCreateEntry(t *testing.T) {
	_, account, err := createRandomAccount()
	require.NoError(t, err)

	arg, entry, err := createRandomEntry(t, account)
	require.NoError(t, err)

	require.NotEmpty(t, entry)
	require.Equal(t, arg.AccountID, entry.AccountID)
	require.Equal(t, arg.Amount, entry.Amount)
	require.NotZero(t, entry.ID)
	require.NotZero(t, entry.CreatedAt)
}

func TestGetEntry(t *testing.T) {
	_, account, err := createRandomAccount()
	require.NoError(t, err)

	_, entry, err := createRandomEntry(t, account)
	require.NoError(t, err)

	entryGet, err := testQueries.GetEntry(context.Background(), entry.ID)
	require.NoError(t, err)
	require.NotEmpty(t, entryGet)

	require.Equal(t, entry.ID, entryGet.ID)
	require.Equal(t, entry.AccountID, entryGet.AccountID)
	require.Equal(t, entry.Amount, entryGet.Amount)
	require.WithinDuration(t, entry.CreatedAt, entryGet.CreatedAt, time.Second)
}

func TestListEntries(t *testing.T) {
	_, account, err := createRandomAccount()
	require.NoError(t, err)

	for i := 0; i < 5; i++ {
		_, _, err := createRandomEntry(t, account)
		require.NoError(t, err)
	}

	arg := ListEntriesParams{Limit: 5, Offset: 0, AccountID: account.ID}

	entries, err := testQueries.ListEntries(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, entries, 5)

	for _, entry := range entries {
		require.NotEmpty(t, entry)
	}
}

func createRandomEntry(t *testing.T, account Account) (Entry, Entry, error) {
	arg := CreateEntryParams{
		AccountID: account.ID,
		Amount:    util.RandomBalance(),
	}

	createdEntry, err := testQueries.CreateEntry(context.Background(), arg)
	return Entry{AccountID: arg.AccountID, Amount: arg.Amount}, createdEntry, err
}
