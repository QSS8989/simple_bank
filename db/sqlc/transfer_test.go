package db

import (
	"context"
	"simple_bank/util"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func createRandomTransfer(t *testing.T, from_account, to_account Account) Transfer {
	arg := CreateTransferParams{
		FromAccountID: from_account.ID,
		ToAccountID:   to_account.ID,
		Amount:        util.RandomMoney(),
	}

	result, err := testQueries.CreateTransfer(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, result)

	require.Equal(t, arg.FromAccountID, result.FromAccountID)
	require.Equal(t, arg.ToAccountID, result.ToAccountID)
	require.Equal(t, arg.Amount, result.Amount)

	require.NotZero(t, result.ID)
	require.NotZero(t, result.CreateAt)

	return result
}

func TestCreateTransfer(t *testing.T) {
	first_account := createRandomAccount(t)
	second_account := createRandomAccount(t)
	createRandomTransfer(t, first_account, second_account)
}

func TestGetTransfer(t *testing.T) {
	first_account := createRandomAccount(t)
	second_account := createRandomAccount(t)
	first_transfer := createRandomTransfer(t, first_account, second_account)

	second_transfer, err := testQueries.GetTransfer(context.Background(), first_transfer.ID)
	require.NoError(t, err)
	require.NotEmpty(t, second_transfer)

	require.Equal(t, first_transfer.ID, second_transfer.ID)
	require.Equal(t, first_transfer.FromAccountID, second_transfer.FromAccountID)
	require.Equal(t, first_transfer.ToAccountID, second_transfer.ToAccountID)
	require.Equal(t, first_transfer.Amount, second_transfer.Amount)
	require.WithinDuration(t, first_transfer.CreateAt, second_transfer.CreateAt, time.Second)
}

func TestListTransfer(t *testing.T) {
	from_account := createRandomAccount(t)
	to_account := createRandomAccount(t)

	for i := 0; i < 10; i++ {
		createRandomTransfer(t, from_account, to_account)
		createRandomTransfer(t, to_account, from_account)
	}

	arg := ListTransfersParams{
		FromAccountID: from_account.ID,
		ToAccountID:   to_account.ID,
		Limit:         5,
		Offset:        5,
	}

	transfers, err := testQueries.ListTransfers(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, transfers, 5)

	for _, transfer := range transfers {
		require.NotEmpty(t, transfer)
		require.True(t, transfer.FromAccountID == from_account.ID || transfer.ToAccountID == from_account.ID)
	}
}
