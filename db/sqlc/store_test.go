package db

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTransferTx(t *testing.T) {
	store := NewStore(testDB)

	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)
	fmt.Println(">>>>Before", account1.Balance, account2.Balance)

	// run n concurrent transfer transaction

	n := 5
	amount := int64(10)
	errs := make(chan error)
	results := make(chan TransferTxResult)
	for i := 0; i < n; i++ {
		//this go routine is running in a  different goroutine than our
		//test functikon is running on , so there is no gurantee taht
		// it will stop the whole condn if a test condition is not
		// satisfied that's why we create a channel

		//solving db deadlock
		// step 1: adding logs

		// txName := fmt.Sprintf("tx %d", i+1) //deadloack solved
		go func() {
			//step 2: db deadlock
			// ctx := context.WithValue(context.Background(), txKey, txName) //deadloack solved
			result, err := store.TransferTx(context.Background(), CreateTransferParams{
				FromAccountID: account1.ID,
				ToAccountID:   account2.ID,
				Amount:        amount,
			})
			errs <- err
			results <- result
		}()
	}

	existed := make(map[int]bool)
	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err, "1")
		result := <-results
		require.NotEmpty(t, result, "2")
		require.NotEmpty(t, result.Transfer, "3")
		require.Equal(t, account1.ID, result.Transfer.FromAccountID, "4")
		require.Equal(t, account2.ID, result.Transfer.ToAccountID, "5")
		require.Equal(t, amount, result.Transfer.Amount, "6")
		require.NotZero(t, result.Transfer.ID, "7")
		require.NotZero(t, result.Transfer.CreatedAt, "8")

		_, err = store.GetTransfer(context.Background(), result.Transfer.ID)
		require.NoError(t, err, "9")

		// check entry
		fromEntry := result.FromEntry
		require.NotEmpty(t, fromEntry)
		require.Equal(t, account1.ID, fromEntry.AccountID)
		require.Equal(t, -amount, fromEntry.Amount)
		require.NotZero(t, fromEntry.ID)
		require.NotZero(t, fromEntry.CreatedAt)

		_, err = store.GetEntry(context.Background(), fromEntry.ID)
		require.NoError(t, err)

		toEntry := result.ToEntry
		require.NotEmpty(t, toEntry)
		require.Equal(t, account2.ID, toEntry.AccountID)
		require.Equal(t, amount, toEntry.Amount)
		require.NotZero(t, toEntry.ID)
		require.NotZero(t, toEntry.CreatedAt)

		_, err = store.GetEntry(context.Background(), toEntry.ID)
		require.NoError(t, err)

		// TODO: check balance

		//TDD-test driven development
		fromAccount := result.FromAccount
		require.NotEmpty(t, fromAccount)
		require.Equal(t, account1.ID, fromAccount.ID)

		toAccount := result.ToAccount
		require.NotEmpty(t, toAccount)
		require.Equal(t, account2.ID, toAccount.ID)
		fmt.Println(">>>>tx1", fromAccount.Balance, toAccount.Balance)

		diff1 := account1.Balance - fromAccount.Balance
		diff2 := toAccount.Balance - account2.Balance

		require.Equal(t, diff1, diff2)
		require.True(t, diff1 >= 0)
		require.True(t, diff1%amount == 0)

		k := int(diff1 / amount)
		fmt.Println("acc", k)
		require.True(t, k >= 1 && k <= n)
		require.NotContains(t, existed, k)
		existed[k] = true
	}

	//check updated balance
	updateAccount1, err := testQueries.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)

	updateAccount2, err := testQueries.GetAccount(context.Background(), account2.ID)
	require.NoError(t, err)
	fmt.Println(">>>>Updated", updateAccount1.Balance, updateAccount2.Balance)

	require.Equal(t, account1.Balance-int64(n)*amount, updateAccount1.Balance)
	require.Equal(t, account2.Balance+int64(n)*amount, updateAccount2.Balance)

}

func TestTransferTxDeadlock(t *testing.T) {
	store := NewStore(testDB)

	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)
	fmt.Println(">>>>Before", account1.Balance, account2.Balance)

	// run n concurrent transfer transaction

	n := 10
	amount := int64(10)
	errs := make(chan error)
	for i := 0; i < n; i++ {

		fromAccountId := account1.ID
		toAccountId := account2.ID

		if i%2 == 1 {
			toAccountId, fromAccountId = account1.ID, account2.ID
		}
		go func() {
			//step 2: db deadlock
			// ctx := context.WithValue(context.Background(), txKey, txName) //deadloack solved
			_, err := store.TransferTx(context.Background(), CreateTransferParams{
				FromAccountID: fromAccountId,
				ToAccountID:   toAccountId,
				Amount:        amount,
			})
			errs <- err
			// results <- result
		}()
	}

	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)
	}

	//check updated balance
	updateAccount1, err := testQueries.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)

	updateAccount2, err := testQueries.GetAccount(context.Background(), account2.ID)
	require.NoError(t, err)
	fmt.Println(">>>>Updated", updateAccount1.Balance, updateAccount2.Balance)

	require.Equal(t, account1.Balance, updateAccount1.Balance)
	require.Equal(t, account2.Balance, updateAccount2.Balance)

}
