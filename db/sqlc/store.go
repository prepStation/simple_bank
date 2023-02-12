package db

import (
	"context"
	"database/sql"
	"fmt"
)

//Db transaction

// Defines all the functionalities needed to execute db transaction
// and Queries
// Because 	the queries struct only define functions
// that can retrieve or insert data on a single table
type Store struct {
	*Queries
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{
		db:      db,
		Queries: New(db),
	}
}

// execTx executes a function within a db transaction
func (store *Store) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		fmt.Println("err1")
		return err
	}

	q := New(tx)
	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			fmt.Println("err2")

			return fmt.Errorf("tx error %v, rbErr %v", err, rbErr)
		}
		fmt.Println("err3")
		return err
	}
	return tx.Commit()
}

// TransferTxParams defines the input for Transfer transaction
// type TransferTxParams struct {
// 	FromAccountID int64 `json:"from_account_id"`
// 	ToAccountID   int64 `json:"to_account_id"`
// 	Amount        int64 `json:"amount"`
// }

// TransferTxResult defines the result of transfer transaction
type TransferTxResult struct {
	Transfer    Transfer `json:"transfer"`
	FromAccount Account  `json:"from_account_id"`
	ToAccount   Account  `json:"to_account_id"`
	FromEntry   Entry    `json:"from_entry"`
	ToEntry     Entry    `json:"to_entry"`
}

//step3: for debugging db deadlock
//Creating our own custom type
// we have to use this key to get the transaction name  from the
//input context of the TransferTx function

// var txKey = struct{}{} deadloack solved

// TransferTx performss money transfer transaction between two accounts
// It Creates transfer record, add account entries, update account Balanace 's within a single
// db transaction
func (store *Store) TransferTx(ctx context.Context, arg CreateTransferParams) (TransferTxResult, error) {
	var result TransferTxResult

	err := store.execTx(ctx, func(q *Queries) error {

		var err error

		// txName := ctx.Value(txKey) deadloack solved

		// fmt.Println(txName, "CreateTransfer") deadloack solved
		result.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams{
			FromAccountID: arg.FromAccountID,
			ToAccountID:   arg.ToAccountID,
			Amount:        arg.Amount,
		})
		if err != nil {
			return err
		}
		// fmt.Println(txName, "CreateEntry1") deadloack solved

		result.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.FromAccountID,
			Amount:    -arg.Amount,
		})
		if err != nil {
			return err
		}
		// fmt.Println(txName, "CreateEntry2") deadloack solved

		result.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.ToAccountID,
			Amount:    arg.Amount,
		})
		if err != nil {
			return err
		}
		//TODO : Update balance accounts
		//How to change account's balance
		//get account from the database
		//add or substract some amount of money from its balance
		//and update it back to the database
		// fmt.Println(txName, "GetAccountForUpdate1") deadloack solved

		//remove this piece code as we are going to converted
		//this to a single query --last

		// account1, err := q.GetAccountForUpdate(ctx, arg.FromAccountID)
		// if err != nil {
		// 	return err
		// }
		// fmt.Println(txName, "UpdateAccount1") deadloack solved

		// result.FromAccount, err = q.UpdateAccount(ctx, UpdateAccountParams{
		// 	ID:      arg.FromAccountID,
		// 	Balance: account1.Balance - arg.Amount,
		// })

		//avoid reverse transaction deadlock

		//lecture 8
		if arg.FromAccountID < arg.ToAccountID {
			//lecture 8
			result.FromAccount, result.ToAccount, err = addMoney(
				ctx, q, arg.FromAccountID,
				-arg.Amount, arg.ToAccountID, arg.Amount,
			)

			if err != nil {
				return err
			}
		} else {
			result.ToAccount, result.FromAccount, err = addMoney(
				ctx, q, arg.ToAccountID,
				arg.Amount, arg.FromAccountID, -arg.Amount,
			)
			if err != nil {
				return err
			}
		}
		// last
		// result.FromAccount, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
		// 	ID:     arg.FromAccountID,
		// 	Amount: -arg.Amount,
		// })
		// if err != nil {
		// 	return err
		// }
		// fmt.Println(txName, "GetAccountForUpdate2") deadloack solved

		//last
		// account2, err := q.GetAccountForUpdate(ctx, arg.ToAccountID)
		// if err != nil {
		// 	return err
		// }
		// fmt.Println(txName, "UpdateAccount2") deadloack solved

		//last
		// result.ToAccount, err = q.UpdateAccount(ctx, UpdateAccountParams{
		// 	ID:      arg.ToAccountID,
		// 	Balance: account2.Balance + arg.Amount,
		// })

		// result.ToAccount, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
		// 	ID:     arg.ToAccountID,
		// 	Amount: arg.Amount,
		// })

		// if err != nil {
		// 	return err
		// }
		return nil
	})

	return result, err
}

// lecture 8
func addMoney(
	ctx context.Context,
	q *Queries,
	accountID1, amount1, accountID2, amount2 int64,
) (account1 Account, account2 Account, err error) {
	account1, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
		Amount: amount1,
		ID:     accountID1,
	})

	if err != nil {
		return
	}

	account2, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
		Amount: amount2,
		ID:     accountID2,
	})

	return
}
