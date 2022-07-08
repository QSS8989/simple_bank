package db

import (
	"context"
	"simple_bank/util"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func createRandomEntry(t *testing.T, account Account) Entry {
	arg := CreateEntryParams{
		AccountID: account.ID,
		Amount:    util.RandomMoney(),
	}

	result, err := testQueries.CreateEntry(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, result)

	require.Equal(t, arg.AccountID, result.AccountID)
	require.Equal(t, arg.Amount, result.Amount)

	require.NotZero(t, result.ID)
	require.NotZero(t, result.CreateAt)

	return result
}

func TestCreateEntry(t *testing.T) {
	account := createRandomAccount(t)
	createRandomEntry(t, account)
}

func TestGetEntry(t *testing.T) {
	account := createRandomAccount(t)
	entry := createRandomEntry(t, account)
	result, err := testQueries.GetEntry(context.Background(), entry.ID)
	require.NoError(t, err)
	require.NotEmpty(t, result)

	require.Equal(t, entry.ID, result.ID)
	require.Equal(t, entry.AccountID, result.AccountID)
	require.Equal(t, entry.Amount, result.Amount)
	require.WithinDuration(t, entry.CreateAt, result.CreateAt, time.Second)
}

func TestListEntries(t *testing.T) {
	account := createRandomAccount(t)
	for i := 0; i < 10; i++ {
		createRandomEntry(t, account)
	}

	arg := ListEntriesParams{
		AccountID: account.ID,
		Limit:     5,
		Offset:    5,
	}

	result, err := testQueries.ListEntries(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, result, 5)

	for _, entry := range result {
		require.NotEmpty(t, entry)
		require.Equal(t, arg.AccountID, entry.AccountID)
	}
}
