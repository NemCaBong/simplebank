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
	// trong việc chuyển khoản ta có 5 bước
	// ta sử dụng n concurrent để handle
	concurrencyStep := 5
	amountTransfer := int64(10)
	existed := make(map[int]bool)
	errsChan := make(chan error)
	resultsChan := make(chan TransferTxResult)
	for i := 0; i < concurrencyStep; i++ {
		// start a new routine
		go func() {
			// start the transfer transaction
			transferResult, err := store.TransferTx(context.Background(), TransferTxParams{
				FromAccountID: account1.ID,
				ToAccountID:   account2.ID,
				Amount:        amountTransfer,
			})
			// Vấn đề ở đây chính là chúng ta phải chạy
			// các bước trong 1 transactions 1 cách concurrent vs nhau
			// cần sử dụng go
			// nhưng mà vì thế chúng ta sẽ tạo ra 1 concurrency chạy cùng với main concurrency
			// không có gì đảm bảo main chạy chậm hơn chúng ta
			// vì thế chúng ta phải để kết quả testify require trong 1 channel
			// và đẩy nó ra ngoài func main ở ngoài.
			// để đảm bảo kết quả đủ đk cho mọi require
			errsChan <- err
			resultsChan <- transferResult
		}()
	}
	// check the result
	for i := 0; i < concurrencyStep; i++ {
		err := <-errsChan
		require.NoError(t, err)

		result := <-resultsChan
		require.NotEmpty(t, result)

		// check the transfer
		transfer := result.Transfer
		require.NotEmpty(t, transfer)
		require.Equal(t, transfer.FromAccountID, account1.ID)
		require.Equal(t, transfer.ToAccountID, account2.ID)
		require.Equal(t, transfer.Amount, amountTransfer)
		require.NotZero(t, transfer.CreatedAt)
		require.NotZero(t, transfer.ID)

		// recheck if the transfer exist or not
		_, err = store.GetTransfer(context.Background(), transfer.ID)
		require.NoError(t, err)

		// check the fromEntry
		fromEntry := result.FromEntry
		require.NotEmpty(t, fromEntry)
		require.Equal(t, fromEntry.Amount, -amountTransfer)
		require.Equal(t, fromEntry.AccountID, account1.ID)
		require.NotZero(t, fromEntry.ID)
		require.NotZero(t, fromEntry.CreatedAt)

		// check if the fromEntry exist in the db or not
		_, err = store.GetEntry(context.Background(), fromEntry.ID)
		require.NoError(t, err)

		// check the toEntry
		toEntry := result.ToEntry
		require.NotEmpty(t, toEntry)
		require.Equal(t, toEntry.Amount, amountTransfer)
		require.Equal(t, toEntry.AccountID, account2.ID)
		require.NotZero(t, toEntry.ID)
		require.NotZero(t, toEntry.CreatedAt)

		// check if the toEntry exist in the db or not
		_, err = store.GetEntry(context.Background(), toEntry.ID)
		require.NoError(t, err)

		//TODO: check the account balances later

		//*c Check the fromAccount
		fromAccount := result.FromAccount
		require.NotEmpty(t, fromAccount)
		require.Equal(t, fromAccount.ID, account1.ID)
		require.NotZero(t, fromAccount.CreatedAt)

		//* check the toAccount
		toAccount := result.ToAccount
		require.NotEmpty(t, toAccount)
		require.Equal(t, toAccount.ID, account2.ID)
		require.NotZero(t, fromAccount.CreatedAt)

		// * check the Account balance
		diff1 := account1.Balance - fromAccount.Balance // fromAccount has been -10
		diff2 := toAccount.Balance - account2.Balance   // toAccount has been +10
		require.Equal(t, diff1, diff2)
		require.True(t, diff1 > 0)                 // so diff2 > 0 too
		require.True(t, diff1%amountTransfer == 0) // số khác bt phải chia hết cho lượng tiền chuyển

		transferTimes := int(diff1 / amountTransfer)
		require.True(t, transferTimes >= 1 && transferTimes <= concurrencyStep)
		// create a slice existed to record the times of the transferTimes
		require.NotContains(t, existed, transferTimes)
		// if it not the same
		existed[transferTimes] = true

	}
	// check the updated balances
	updatedAccount1, err := store.GetAccountForUpdate(context.Background(), account1.ID)
	require.NoError(t, err)

	updatedAccount2, err := store.GetAccountForUpdate(context.Background(), account2.ID)
	require.NoError(t, err)

	require.Equal(t, account1.Balance-int64(concurrencyStep)*amountTransfer, updatedAccount1.Balance)
	require.Equal(t, account2.Balance+int64(concurrencyStep)*amountTransfer, updatedAccount2.Balance)

}

func TestTransferTxDeadlock(t *testing.T) {
	store := NewStore(testDB)
	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)
	fmt.Println(">>> before:", account1.Balance, account2.Balance)

	concurrencyStep := 10
	amountTransfer := int64(10)
	errsChan := make(chan error)
	for i := 0; i < concurrencyStep; i++ {
		// start a new routine
		fromAccountID := account1.ID
		toAccountId := account2.ID
		if i%2 == 1 {
			fromAccountID = account2.ID
			toAccountId = account1.ID
		}
		go func() {
			// start the transfer transaction
			_, err := store.TransferTx(context.Background(), TransferTxParams{
				FromAccountID: fromAccountID,
				ToAccountID:   toAccountId,
				Amount:        amountTransfer,
			})
			errsChan <- err
		}()
	}
	// check the result
	for i := 0; i < concurrencyStep; i++ {
		err := <-errsChan
		require.NoError(t, err)
		// check the updated balances
	}
	// check outside because each 2 steps the balance can be balanced
	updatedAccount1, err := testQueries.GetAccountForUpdate(context.Background(), account1.ID)
	require.NoError(t, err)

	updatedAccount2, err := testQueries.GetAccountForUpdate(context.Background(), account2.ID)
	require.NoError(t, err)
	fmt.Println(">>>> After:", updatedAccount1.Balance, updatedAccount2.Balance)

	require.Equal(t, account1.Balance, updatedAccount1.Balance)
	require.Equal(t, account2.Balance, updatedAccount2.Balance)
}
