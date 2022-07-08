package db

import (
	"context"
	"database/sql"
	"simple_bank/util"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func createRandomAccount(t *testing.T) Account {
	arg := CreateAccountParams{
		Owner:    util.RandOwner(),
		Blance:   util.RandomMoney(),
		Currency: util.RandCurrency(),
	}

	account, err := testQueries.CreateAccount(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, account)
	require.Equal(t, account.Owner, arg.Owner)
	require.Equal(t, account.Blance, arg.Blance)
	require.Equal(t, account.Currency, arg.Currency)
	require.NotZero(t, account.ID)
	require.NotZero(t, account.CreateAt)
	return account
}
func TestCreateAccount(t *testing.T) {
	createRandomAccount(t)
}

func TestGetAccount(t *testing.T) {
	account_c := createRandomAccount(t)
	account_g, err := testQueries.GetAccount(context.Background(), account_c.ID)
	require.NoError(t, err)
	require.NotEmpty(t, account_g)

	require.Equal(t, account_c.ID, account_g.ID)
	require.Equal(t, account_c.Owner, account_g.Owner)
	require.Equal(t, account_c.Blance, account_g.Blance)
	require.Equal(t, account_c.Currency, account_g.Currency)
	require.WithinDuration(t, account_c.CreateAt, account_g.CreateAt, time.Second)
}

func TestUpdateAccount(t *testing.T) {
	account_c := createRandomAccount(t)
	arg := UpdateAccountParams{
		ID:     account_c.ID,
		Blance: util.RandomMoney(),
	}
	account_u, err := testQueries.UpdateAccount(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, account_u)

	require.Equal(t, account_c.ID, account_u.ID)
	require.Equal(t, account_c.Owner, account_u.Owner)
	require.Equal(t, arg.Blance, account_u.Blance)
	require.Equal(t, account_c.Currency, account_u.Currency)
	require.WithinDuration(t, account_c.CreateAt, account_u.CreateAt, time.Second)

}

func TestDeleteAccount(t *testing.T) {
	account_c := createRandomAccount(t)
	err := testQueries.DeleteAccount(context.Background(), account_c.ID)
	require.NoError(t, err)

	account_g, err := testQueries.GetAccount(context.Background(), account_c.ID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, account_g)

}

func TestListAccount(t *testing.T) {
	for i := 0; i < 10; i++ {
		createRandomAccount(t)
	}

	arg := ListAccountsParams{
		Limit:  5,
		Offset: 5,
	}

	account_l, err := testQueries.ListAccounts(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, account_l, 5)

	for _, account := range account_l {
		require.NotEmpty(t, account)
	}

}
