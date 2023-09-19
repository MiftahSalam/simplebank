package db

import (
	"context"
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

	accountGet, err := testStore.GetAccount(context.Background(), account.ID)
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
	accountUpdated, err := testStore.UpdateAccountBalance(context.Background(), arg)
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

	err = testStore.DeleteAccount(context.Background(), account.ID)
	require.NoError(t, err)

	accountGet, err := testStore.GetAccount(context.Background(), account.ID)
	require.Error(t, err)
	require.EqualError(t, err, ErrRecordNotFound.Error())
	require.Empty(t, accountGet)

}

func TestListAccount(t *testing.T) {
	var lastAccount Account
	for i := 0; i < 5; i++ {
		var err error
		_, lastAccount, err = createRandomAccount()
		require.NoError(t, err)
	}

	arg := ListAccountParams{Owner: lastAccount.Owner, Limit: 5, Offset: 0}

	accounts, err := testStore.ListAccount(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, accounts)

	for _, account := range accounts {
		require.NotEmpty(t, account)
		require.Equal(t, account.Owner, lastAccount.Owner)
	}
}

func createRandomAccount() (Account, Account, error) {
	_, user, _ := createRandomUser()
	arg := CreateAccountParams{
		Owner:    user.Username,
		Balance:  util.RandomBalance(),
		Currency: util.RandomCurrency(),
	}

	createdAccount, err := testStore.CreateAccount(context.Background(), arg)
	return Account{Owner: arg.Owner, Balance: arg.Balance, Currency: arg.Currency}, createdAccount, err
}
