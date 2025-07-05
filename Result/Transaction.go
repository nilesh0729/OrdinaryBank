package main

import (
	"context"
	"database/sql"
	"fmt"
)

type Store struct{
	*Queries
	db *sql.DB
}

func NewStore(db *sql.DB) *Store{
	return &Store{
		db : db,
		Queries: New(db),
	}
}

func (store *Store) exexTx(ctx context.Context, fn func(*Queries)error)error{
	tx, err := store.db.BeginTx(ctx,nil)
	if err != nil {
		return err
	}
	q := New(tx)
	err = fn(q)

	if err != nil{
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("rollback Err: %v, TxErr: %v", rbErr, err)
		}
		return err
	}
	return tx.Commit()

}

type TransferTxParams struct {
	FromAccountID int64 
	ToAccountID int64
	Amount int64
}

type TransferTxResult struct{
	Transfer Transfer
	FromAccountID Account
	ToAccountID Account
	FromEntry Entry
	ToEntry Entry
}

func (store *Store) TransferTx(ctx context.Context, arg TransferTxParams)(TransferTxResult, error){
	var result TransferTxResult
	err := store.exexTx(ctx, func(q *Queries)error{
		var err error
		result.Transfer, err = q.CreateTransfers(ctx, CreateTransfersParams{
			FromAccountID: arg.FromAccountID,
			ToAccountID: arg.ToAccountID,
			Amount: arg.Amount,
		})
		if err != nil {
			return err
		}
		result.FromEntry, err = q.CreateEntries(ctx, CreateEntriesParams{
			AccountID: arg.FromAccountID,
			Amount: -arg.Amount,
		})
		if err != nil{
			return err
		}

		result.ToEntry , err = q.CreateEntries(ctx, CreateEntriesParams{
			AccountID: arg.ToAccountID,
			Amount: arg.Amount,
		})

		if err != nil{
			return err
		}

		// TODO := THE UPDATE BALANCE(ADD AMOUNT IN ACCOUNT1 & SUBTRACT AMOUNT
		//  FROM ACCOUNT2) CODE WILL BE ADDED SOON 

		return err
	})

	return result , err
}

