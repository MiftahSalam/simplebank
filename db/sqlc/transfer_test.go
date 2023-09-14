package db

import (
	"context"
	"simplebank/util"
	"testing"
	"time"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/require"
)

func TestCreateTransfer(t *testing.T) {
	_, accountFrom, err := createRandomAccount()
	require.NoError(t, err)
	_, accountTo, err := createRandomAccount()
	require.NoError(t, err)

	arg, tf, err := createRandomTransfer(t, accountTo, accountFrom)
	require.NoError(t, err)

	require.NotEmpty(t, tf)
	require.Equal(t, arg.ToAccountID, tf.ToAccountID)
	require.Equal(t, arg.FromAccountID, tf.FromAccountID)
	require.Equal(t, arg.Amount, tf.Amount)
	require.NotZero(t, tf.ID)
	require.NotZero(t, tf.CreatedAt)
}

func TestGetTransfer(t *testing.T) {
	_, accountFrom, err := createRandomAccount()
	require.NoError(t, err)
	_, accountTo, err := createRandomAccount()
	require.NoError(t, err)

	_, tf, err := createRandomTransfer(t, accountTo, accountFrom)
	require.NoError(t, err)

	tfGet, err := testQueries.GetTransfer(context.Background(), tf.ID)
	require.NoError(t, err)
	require.NotEmpty(t, tfGet)

	require.Equal(t, tf.ID, tfGet.ID)
	require.Equal(t, tf.FromAccountID, tfGet.FromAccountID)
	require.Equal(t, tf.ToAccountID, tfGet.ToAccountID)
	require.Equal(t, tf.Amount, tfGet.Amount)
	require.WithinDuration(t, tf.CreatedAt, tfGet.CreatedAt, time.Second)
}

func TestListTransfers(t *testing.T) {
	_, accountFrom, err := createRandomAccount()
	require.NoError(t, err)
	_, accountTo, err := createRandomAccount()
	require.NoError(t, err)

	for i := 0; i < 5; i++ {
		_, _, err := createRandomTransfer(t, accountFrom, accountTo)
		require.NoError(t, err)
	}

	arg := ListTransfersParams{Limit: 5, Offset: 0, FromAccountID: accountFrom.ID, ToAccountID: accountTo.ID}

	tfs, err := testQueries.ListTransfers(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, tfs, 5)

	for _, tf := range tfs {
		require.NotEmpty(t, tf)
	}
}

func createRandomTransfer(t *testing.T, accountFrom, accountTo Account) (Transfer, Transfer, error) {
	arg := CreateTransferParams{
		FromAccountID: accountFrom.ID,
		ToAccountID:   accountTo.ID,
		Amount:        util.RandomBalance(),
	}

	createdTf, err := testQueries.CreateTransfer(context.Background(), arg)
	return Transfer{FromAccountID: arg.FromAccountID, ToAccountID: accountTo.ID, Amount: arg.Amount}, createdTf, err
}
