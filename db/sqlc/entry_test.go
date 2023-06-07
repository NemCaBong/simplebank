package db

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/techschool/simplebank/db/util"
)

func createRandomEntry(t *testing.T, account Account) Entry {
	arg := CreateEntryParams{
		AccountID: account.ID,
		Amount:    util.RandomMoney(),
	}
	entry, err := testQueries.CreateEntry(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, entry)
	require.Equal(t, arg.Amount, entry.Amount)
	require.Equal(t, arg.AccountID, entry.AccountID)
	require.NotZero(t, entry.ID)
	require.NotZero(t, entry.CreatedAt)

	return entry
}
func TestCreateEntry(t *testing.T) {
	account := createRandomAccount(t)
	createRandomEntry(t, account)
}

func TestGetEntry(t *testing.T) {
	account := createRandomAccount(t)
	entry1 := createRandomEntry(t, account)
	entry2, err := testQueries.GetEntry(context.Background(), entry1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, entry2)
	require.Equal(t, entry1.ID, entry2.ID)
	require.Equal(t, entry1.AccountID, entry2.AccountID)
	require.Equal(t, entry1.Amount, entry2.Amount)
	require.WithinDuration(t, entry1.CreatedAt, entry2.CreatedAt, time.Second)
}

func TestListEntry(t *testing.T) {
	account := createRandomAccount(t)
	// using that account id to create 10 entries.
	for i := 0; i < 10; i++ {
		createRandomEntry(t, account)
	}
	// account.ID = entry.AccountID mà thôi
	// bởi vì mình chỉ dùng account đó
	arg := ListEntriesParams{
		AccountID: account.ID,
		Offset:    5,
		Limit:     5,
	}
	listEntries, err := testQueries.ListEntries(context.Background(), arg)
	require.NoError(t, err)
	// đảm bảo có 5 entries 1 lượt
	require.Len(t, listEntries, 5)
	for _, entry := range listEntries {
		// check xem có empty hay không
		// phải để ở trc
		require.NotEmpty(t, entry)
		require.Equal(t, arg.AccountID, entry.AccountID)
	}
}
