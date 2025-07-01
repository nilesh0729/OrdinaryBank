package main

import (
	"context"
	"learning/util"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func CreateRandomTransfers(t *testing.T, account1, account2 Account)Transfer {
	arg := CreateTransfersParams{
		FromAccountID: account1.ID,
		ToAccountID: account2.ID,

		Amount: util.RandomBalance(),
	}

	Transfers, err := testQueries.CreateTransfers(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, Transfers)

	require.Equal(t,Transfers.FromAccountID, arg.FromAccountID)
	require.Equal(t, Transfers.ToAccountID, arg.ToAccountID)
	require.Equal(t, Transfers.Amount, arg.Amount)
    
	require.NotZero(t,Transfers.ID)
	require.NotZero(t, Transfers.CreatedAt)

	return Transfers

}

func TestCreateTransfers(t *testing.T){
	account1 := CreateRandomAccount(t)
	account2 := CreateRandomAccount(t)

	CreateRandomTransfers(t,account1,account2)
}

func TestGetTransfers(t *testing.T){
	account1 := CreateRandomAccount(t)
	account2 := CreateRandomAccount(t)

	transfer1 := CreateRandomTransfers(t, account1,account2)

	transfer2, err:= testQueries.GetTransfers(context.Background(), transfer1.ID)

	require.NoError(t, err)
	require.NotEmpty(t, transfer2)

	require.Equal(t, transfer1.ID, transfer2.ID)
	require.Equal(t, transfer1.FromAccountID, transfer2.FromAccountID)
	require.Equal(t, transfer1.ToAccountID, transfer2.ToAccountID)
	require.Equal(t,transfer1.Amount, transfer2.Amount)

	require.WithinDuration(t, transfer1.CreatedAt, transfer2.CreatedAt, time.Second)
}

func TestListTransfers(t *testing.T) {
	account1 := CreateRandomAccount(t)
	account2 := CreateRandomAccount(t)

	// Create some transfers in both directions
	for i := 0; i < 10; i++ {
		CreateRandomTransfers(t, account1, account2)
		CreateRandomTransfers(t, account2, account1)
	}

	arg := ListTransfersParams{
		FromAccountID: account1.ID,
		ToAccountID:   account1.ID,
		Limit:         5,
		Offset:        5, 
	}

	transfers, err := testQueries.ListTransfers(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, transfers, 5)

	for _, transfer := range transfers {
		require.NotEmpty(t, transfer)
		
		require.True(t, transfer.FromAccountID == account1.ID || transfer.ToAccountID == account1.ID)
	}
}

