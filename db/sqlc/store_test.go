package db

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTransferTx(t *testing.T) {
	store := NewStore(testDB)
	account_from := createRandomAccount(t)
	account_to := createRandomAccount(t)
	n := 5
	amount := int64(10)

	errs := make(chan error)
	results := make(chan TransferTxResult)

	for i := 0; i < n; i++ {

		go func() {

			result, err := store.TransferTx(context.Background(), TransferTxParams{
				FromAccountID: account_from.ID,
				ToAccountID:   account_to.ID,
				Amount:        amount,
			})

			errs <- err
			results <- result

		}()
	}

	// check results
	existed := make(map[int]bool)
	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)

		result := <-results
		require.NotEmpty(t, result)

		//check transfer
		transfer := result.Transfer
		require.NotEmpty(t, transfer)
		require.Equal(t, account_from.ID, transfer.FromAccountID)
		require.Equal(t, account_to.ID, transfer.ToAccountID)
		require.Equal(t, amount, transfer.Amount)
		require.NotZero(t, transfer.ID)
		require.NotZero(t, transfer.CreateAt)

		_, err = store.GetTransfer(context.Background(), transfer.ID)
		require.NoError(t, err)

		// check entry

		fromEntry := result.FromEntry
		require.NotEmpty(t, fromEntry)
		require.Equal(t, account_from.ID, fromEntry.AccountID)
		require.Equal(t, -amount, fromEntry.Amount)
		require.NotZero(t, fromEntry.ID)
		_, err = store.GetEntry(context.Background(), fromEntry.ID)
		require.NoError(t, err)

		toEntry := result.ToEntry
		require.NotEmpty(t, toEntry)
		require.Equal(t, account_to.ID, toEntry.AccountID)
		require.Equal(t, amount, toEntry.Amount)
		require.NotZero(t, toEntry.ID)
		_, err = store.GetEntry(context.Background(), toEntry.ID)
		require.NoError(t, err)

		// check accounts

		fromAccount := result.FromAccount
		require.NotEmpty(t, fromAccount)
		require.Equal(t, account_from.ID, fromAccount.ID)

		toAccount := result.ToAccount

		require.NotEmpty(t, toAccount)
		require.Equal(t, account_to.ID, toAccount.ID)

		// 	check account's balance
		fmt.Println(">> tx:", fromAccount.Balance, toAccount.Balance)

		from_amount := account_from.Balance - fromAccount.Balance
		to_amount := toAccount.Balance - account_to.Balance
		require.Equal(t, from_amount, to_amount)
		require.True(t, from_amount > 0)
		require.True(t, from_amount%amount == 0)

		k := int(from_amount / amount)    // 最大转账次数
		require.True(t, k >= 1 && k <= n) // 判断转账金额是否大于转出账户金额
		require.NotContains(t, existed, k)
		existed[k] = true

	}
	fmt.Println(">> after:", account_from.Balance, account_to.Balance)
	// check the final updated balances
	updatedAccount_from, err := testQueries.GetAccount(context.Background(), account_from.ID)
	require.NoError(t, err)
	require.NotEmpty(t, updatedAccount_from)
	updatedAccount_to, err := testQueries.GetAccount(context.Background(), account_to.ID)
	require.NoError(t, err)
	require.NotEmpty(t, updatedAccount_to)
	require.Equal(t, account_from.Balance-int64(n)*amount, updatedAccount_from.Balance)
	require.Equal(t, account_to.Balance+int64(n)*amount, updatedAccount_to.Balance)
}

func TestDeadLock(t *testing.T) {
	store := NewStore(testDB)
	account_from := createRandomAccount(t)
	account_to := createRandomAccount(t)
	n := 10
	amount := int64(10)

	errs := make(chan error)

	for i := 0; i < n; i++ {
		fromAccountID := account_from.ID
		toAccountID := account_to.ID

		if i%2 == 1 {
			fromAccountID = account_to.ID
			toAccountID = account_from.ID

		}
		go func() {
			_, err := store.TransferTx(context.Background(), TransferTxParams{
				FromAccountID: fromAccountID,
				ToAccountID:   toAccountID,
				Amount:        amount,
			})
			errs <- err

		}()
	}
	fmt.Println(">> after:", account_from.Balance, account_to.Balance)
	updatedAccount_from, err := testQueries.GetAccount(context.Background(), account_from.ID)
	require.NoError(t, err)
	require.NotEmpty(t, updatedAccount_from)
	updatedAccount_to, err := testQueries.GetAccount(context.Background(), account_to.ID)
	require.NoError(t, err)
	require.NotEmpty(t, updatedAccount_to)
	require.Equal(t, account_from.Balance, updatedAccount_from.Balance)
	require.Equal(t, account_to.Balance, updatedAccount_to.Balance)
}
