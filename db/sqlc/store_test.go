package db

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTransferTx(t *testing.T) {
	_, account1, err := createRandomAccount()
	require.NoError(t, err)

	_, account2, err := createRandomAccount()
	require.NoError(t, err)

	fmt.Println(">> before: ", account1.Balance, account2.Balance)

	errsChan := make(chan error)
	resultsChan := make(chan TransferTxResult)

	n := 5
	amount := int64(10)
	for i := 0; i < n; i++ {
		go func() {
			result, err := testStore.TransferTx(context.Background(), CreateTransferParams{
				FromAccountID: account1.ID,
				ToAccountID:   account2.ID,
				Amount:        amount,
			})

			errsChan <- err
			resultsChan <- result
		}()
	}

	existed := make(map[int]bool)
	for i := 0; i < n; i++ {
		err := <-errsChan
		require.NoError(t, err)

		result := <-resultsChan

		require.NotEmpty(t, result)
		require.NotEmpty(t, result.Transfer)
		require.Equal(t, account1.ID, result.Transfer.FromAccountID)
		require.Equal(t, account2.ID, result.Transfer.ToAccountID)
		require.Equal(t, amount, result.Transfer.Amount)
		require.NotZero(t, result.Transfer.ID)
		require.NotZero(t, result.Transfer.CreatedAt)

		_, err = testStore.GetTransfer(context.Background(), result.Transfer.ID)
		require.NoError(t, err)

		require.NotEmpty(t, result.FromEntry)
		require.Equal(t, account1.ID, result.FromEntry.AccountID)
		require.Equal(t, -amount, result.FromEntry.Amount)
		require.NotZero(t, result.FromEntry.ID)
		require.NotZero(t, result.FromEntry.CreatedAt)

		_, err = testStore.GetEntry(context.Background(), result.FromEntry.ID)
		require.NoError(t, err)

		require.NotEmpty(t, result.ToEntry)
		require.Equal(t, account2.ID, result.ToEntry.AccountID)
		require.Equal(t, amount, result.ToEntry.Amount)
		require.NotZero(t, result.ToEntry.ID)
		require.NotZero(t, result.ToEntry.CreatedAt)

		_, err = testStore.GetEntry(context.Background(), result.ToEntry.ID)
		require.NoError(t, err)

		require.NotEmpty(t, result.FromAccount)
		require.Equal(t, account1.ID, result.FromAccount.ID)

		require.NotEmpty(t, result.ToAccount)
		require.Equal(t, account2.ID, result.ToAccount.ID)

		fmt.Println(">> tx: ", result.FromAccount.Balance, result.ToAccount.Balance)

		diff1 := account1.Balance - result.FromAccount.Balance
		diff2 := result.ToAccount.Balance - account2.Balance

		require.Equal(t, diff1, diff2)
		require.True(t, diff1 > 0)
		require.True(t, diff1%amount == 0)

		k := int(diff1 / amount)

		require.True(t, k >= 1 && k <= n)
		require.NotContains(t, existed, k)

		existed[k] = true
	}

	updateAccount1, err := testStore.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)

	updateAccount2, err := testStore.GetAccount(context.Background(), account2.ID)
	require.NoError(t, err)

	fmt.Println(">> after: ", updateAccount1.Balance, updateAccount2.Balance)

	require.Equal(t, account1.Balance-int64(n)*amount, updateAccount1.Balance)
	require.Equal(t, account2.Balance+int64(n)*amount, updateAccount2.Balance)
}
