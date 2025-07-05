package main

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTransferTx(t *testing.T) {
	store := NewStore(TestDb)

	account1 := CreateRandomAccount(t)
	account2 := CreateRandomAccount(t)

	n := 5
	amount := int64(10)

	errs := make(chan error)
	results := make(chan TransferTxResult)

	for i:=0; i < n; i++{
		go func(){
			result, err := store.TransferTx(context.Background(), TransferTxParams{
				FromAccountID: account1.ID,
				ToAccountID: account2.ID,
				Amount: amount,
			})

			errs <- err
			results <- result

		}()
		
	}
	// CHECK RESULTS

	for i:=0; i<n; i++ {
		err:= <-errs

		require.NoError(t, err)
		
		result := <- results
		require.NotEmpty(t, result)

		//check Transfers
		Transfer := result.Transfer

		require.Equal(t, account1.ID, Transfer.FromAccountID)
		require.Equal(t, account2.ID, Transfer.ToAccountID)

		require.Equal(t, amount, Transfer.Amount)

		require.NotZero(t , Transfer.ID)
		require.NotZero(t, Transfer.CreatedAt)

		_, err = store.GetTransfers(context.Background(),Transfer.ID)
		require.NoError(t, err)

		//Check Entries

		FromEntry := result.FromEntry
		require.NotEmpty(t, FromEntry)

		require.Equal(t,account1.ID, FromEntry.AccountID)
		require.Equal(t, -amount, FromEntry.Amount)

		require.NotZero(t, FromEntry.ID)
		require.NotEmpty(t, FromEntry.CreatedAt)

		_,err = store.GetEntries(context.Background(),FromEntry.ID)
		require.NoError(t, err)

		ToEntry := result.ToEntry
		require.NotEmpty(t, ToEntry)

		require.Equal(t, account2.ID, ToEntry.AccountID)
		require.Equal(t, amount, ToEntry.Amount)

		require.NotZero(t, ToEntry.ID)
		require.NotZero(t, ToEntry.CreatedAt)

		_, err = store.GetAccounts(context.Background(), ToEntry.ID)
		require.NoError(t, err)

	}
	
}
