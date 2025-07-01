package main

import (
	"context"
	"learning/util"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func CreateRandomEntries(t *testing.T, account Account)Entry {
	arg :=  CreateEntriesParams{
	AccountID: account.ID,
	Amount: util.RandomBalance(),    
}
    Entry, err := testQueries.CreateEntries(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, Entry)

	require.Equal(t, arg.AccountID, Entry.AccountID)
	require.Equal(t, arg.Amount, Entry.Amount)

	require.NotZero(t, Entry.ID)
	require.NotZero(t, Entry.CreatedAt)

	return Entry
}

func TestCreateEntries(t *testing.T){
	account := CreateRandomAccount(t)
	CreateRandomEntries(t,account)
}

func TestGetEntries(t *testing.T) {
	account := CreateRandomAccount(t)

	entry1 := CreateRandomEntries(t, account)

	entry2, err := testQueries.GetEntries(context.Background(), entry1.ID)

	require.NoError(t,err)
	require.NotEmpty(t,entry2)

	require.Equal(t, entry1.AccountID, entry2.AccountID)
	require.Equal(t, entry1.Amount, entry2.Amount)
	require.Equal(t, entry1.ID, entry2.ID)

	require.WithinDuration(t, entry1.CreatedAt, entry2.CreatedAt, time.Second)

}
func TestListEntries(t *testing.T){
	account := CreateRandomAccount(t)
	for i := 0; i < 10; i++ {
		CreateRandomEntries(t, account)
	}	
	Offset :=   5

	entries, err := testQueries.ListEntries(context.Background(), int32(Offset))
	require.NoError(t, err)
	

	for _, entry := range entries {
		require.NotEmpty(t, entry)
		require.NotEmpty(t, entry.AccountID)
	}


}
