package db

import (
	"context"
	"database/sql"
	"simplebank/util"
	"testing"
	"time"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/require"
)

func TestCreateAccount(t *testing.T) {
	arg, account, err := createRandomAccount()

	require.NoError(t, err)
	require.NotEmpty(t, account)
	require.Equal(t, arg.Owner, account.Owner)
	require.Equal(t, arg.Balance, account.Balance)
	require.Equal(t, arg.Currency, account.Currency)
	require.NotZero(t, account.ID)
	require.NotZero(t, account.CreatedAt)
}

func TestGetAccount(t *testing.T) {
	_, account, err := createRandomAccount()
	require.NoError(t, err)

	accountGet, err := testQueries.GetAccount(context.Background(), account.ID)
	require.NoError(t, err)
	require.NotEmpty(t, accountGet)

	require.Equal(t, account.ID, accountGet.ID)
	require.Equal(t, account.Owner, accountGet.Owner)
	require.Equal(t, account.Balance, accountGet.Balance)
	require.Equal(t, account.Currency, accountGet.Currency)
	require.WithinDuration(t, account.CreatedAt, accountGet.CreatedAt, time.Second)
}

func TestUpdateAccountBalance(t *testing.T) {
	_, account, err := createRandomAccount()
	require.NoError(t, err)

	arg := UpdateAccountBalanceParams{ID: account.ID, Amount: util.RandomBalance()}
	accountUpdated, err := testQueries.UpdateAccountBalance(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, accountUpdated)

	require.Equal(t, account.ID, accountUpdated.ID)
	require.Equal(t, account.Owner, accountUpdated.Owner)
	require.Equal(t, account.Balance+arg.Amount, accountUpdated.Balance)
	require.Equal(t, account.Currency, accountUpdated.Currency)
	require.WithinDuration(t, account.CreatedAt, accountUpdated.CreatedAt, time.Second)
}

func TestDeleteAccount(t *testing.T) {
	_, account, err := createRandomAccount()
	require.NoError(t, err)

	err = testQueries.DeleteAccount(context.Background(), account.ID)
	require.NoError(t, err)

	accountGet, err := testQueries.GetAccount(context.Background(), account.ID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, accountGet)

}

func TestListAccount(t *testing.T) {
	for i := 0; i < 5; i++ {
		_, _, err := createRandomAccount()
		require.NoError(t, err)
	}

	arg := ListAccountParams{Limit: 5, Offset: 5}

	accounts, err := testQueries.ListAccount(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, accounts, 5)

	for _, account := range accounts {
		require.NotEmpty(t, account)
	}
}

func createRandomAccount() (Account, Account, error) {
	arg := CreateAccountParams{
		Owner:    util.RandomOwner(),
		Balance:  util.RandomBalance(),
		Currency: util.RandomCurrency(),
	}

	createdAccount, err := testQueries.CreateAccount(context.Background(), arg)
	return Account{Owner: arg.Owner, Balance: arg.Balance, Currency: arg.Currency}, createdAccount, err
}
