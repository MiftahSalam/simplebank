package db

import (
	"context"
	"simplebank/util"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/require"
)

func TestCreateUser(t *testing.T) {
	arg, user, err := createRandomUser()

	require.NoError(t, err)
	require.NotEmpty(t, user)
	require.Equal(t, arg.Username, user.Username)
	require.Equal(t, arg.FullName, user.FullName)
	require.Equal(t, arg.Email, user.Email)
	require.NotZero(t, user.CreatedAt)
}

func TestGetUser(t *testing.T) {
	_, user, err := createRandomUser()
	require.NoError(t, err)

	userGet, err := testStore.GetUser(context.Background(), user.Username)
	require.NoError(t, err)
	require.NotEmpty(t, userGet)

	require.Equal(t, user.Username, userGet.Username)
	require.Equal(t, user.FullName, userGet.FullName)
	require.Equal(t, user.HashedPassword, userGet.HashedPassword)
	require.Equal(t, user.Email, userGet.Email)
	require.WithinDuration(t, user.PasswordChangedAt, userGet.PasswordChangedAt, time.Second)
	require.WithinDuration(t, user.CreatedAt, userGet.CreatedAt, time.Second)
}

func TestUpdateUser(t *testing.T) {
	_, user, err := createRandomUser()
	require.NoError(t, err)

	updatedUser, err := testStore.UpdateUser(context.Background(), UpdateUserParams{
		Username: user.Username,
		FullName: pgtype.Text{
			String: util.RandomOwner(),
			Valid:  true,
		},
	})

	require.NoError(t, err)
	require.NotEqual(t, user.FullName, updatedUser.FullName)
	require.Equal(t, user.Username, updatedUser.Username)
	require.Equal(t, user.Email, updatedUser.Email)
}

func createRandomUser() (User, User, error) {
	hashedPassword, _ := util.HashPassword(util.RandomString(6))
	arg := CreateUserParams{
		Username:       util.RandomOwner(),
		FullName:       util.RandomOwner(),
		HashedPassword: hashedPassword,
		Email:          util.RandomEmail(),
	}

	createdUser, err := testStore.CreateUser(context.Background(), arg)
	return User{Username: arg.Username, FullName: arg.FullName, HashedPassword: arg.HashedPassword, Email: arg.Email}, createdUser, err
}
