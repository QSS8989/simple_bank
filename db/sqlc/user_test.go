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

	require.True(t, user.PasswordChangeAt.IsZero())
	require.NotZero(t, user.CreateAt)

	return user
}

func TestCreateUser(t *testing.T) {
	createRandomUser(t)
}

func TestGetUser(t *testing.T) {
	user_c := createRandomUser(t)
	user_g, err := testQueries.GetUser(context.Background(), user_c.Username)
	require.NoError(t, err)
	require.NotEmpty(t, user_g)

	require.Equal(t, user_c.Username, user_g.Username)
	require.Equal(t, user_c.HashedPassword, user_g.HashedPassword)
	require.Equal(t, user_c.FullName, user_g.FullName)
	require.Equal(t, user_c.Email, user_g.Email)
	require.WithinDuration(t, user_c.CreateAt, user_g.CreateAt, time.Second)
	require.WithinDuration(t, user_c.PasswordChangeAt, user_g.PasswordChangeAt, time.Second)

}
