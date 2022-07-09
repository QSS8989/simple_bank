package db

import (
	"context"
	"simple_bank/util"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func createRandomUser(t *testing.T) User {
	hashedPassword, err := util.HashPassword(util.RandomString(6))
	require.NoError(t, err)
	require.NotEmpty(t, hashedPassword)

	arg := CreateUserParams{
		Username:       util.RandomOwner(),
		HashedPassword: hashedPassword,
		FullName:       util.RandomOwner(),
		Email:          util.RandomEmail(),
	}

	user, err := testQueries.CreateUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, user)

	require.Equal(t, arg.Username, user.Username)
	require.Equal(t, arg.HashedPassword, user.HashedPassword)
	require.Equal(t, arg.FullName, user.FullName)
	require.Equal(t, arg.Email, user.Email)

	require.True(t, user.PasswordChangedAt.IsZero())
	require.NotZero(t, user.CreatedAt)

	return user
}

func TestCreateUser(t *testing.T) {
	createRandomUser(t)
}

func TestGetUser(t *testing.T) {
	newUser := createRandomUser(t)
	result, err := testQueries.GetUser(context.Background(), newUser.Username)
	require.NoError(t, err)
	require.NotEmpty(t, result)

	require.Equal(t, newUser.Username, result.Username)
	require.Equal(t, newUser.HashedPassword, result.HashedPassword)
	require.Equal(t, newUser.FullName, result.FullName)
	require.Equal(t, newUser.Email, result.Email)
	require.WithinDuration(t, newUser.CreatedAt, result.CreatedAt, time.Second)
	require.WithinDuration(t, newUser.PasswordChangedAt, result.PasswordChangedAt, time.Second)

}
