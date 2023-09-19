package db

import "context"

type TransferTxResult struct {
	Transfer    Transfer `json:"transfer"`
	FromAccount Account  `json:"from_account"`
	ToAccount   Account  `json:"to_account"`
	FromEntry   Entry    `json:"from_entry"`
	ToEntry     Entry    `json:"to_entry"`
}

func (store *SQLStore) TransferTx(ctx context.Context, arg CreateTransferParams) (TransferTxResult, error) {
	var result TransferTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		result.Transfer, err = q.CreateTransfer(ctx, arg)
		if err != nil {
			return err
		}

		result.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.FromAccountID,
			Amount:    -arg.Amount,
		})
		if err != nil {
			return err
		}

		result.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.ToAccountID,
			Amount:    arg.Amount,
		})
		if err != nil {
			return err
		}

		if result.FromAccount.ID < result.Transfer.ID {
			result.FromAccount, result.ToAccount, err = moneyMovement(ctx, q, arg.FromAccountID, -arg.Amount, arg.ToAccountID, arg.Amount)
			if err != nil {
				return err
			}
		} else {
			result.ToAccount, result.FromAccount, err = moneyMovement(ctx, q, arg.ToAccountID, arg.Amount, arg.FromAccountID, -arg.Amount)
			if err != nil {
				return err
			}
			if err != nil {
				return err
			}
		}

		return nil
	})

	return result, err
}

func moneyMovement(ctx context.Context, q *Queries, accountIdFrom, amountFrom, accountIdTo, amountTo int64) (Account, Account, error) {
	accountFrom, err := q.UpdateAccountBalance(ctx, UpdateAccountBalanceParams{
		ID:     accountIdFrom,
		Amount: amountFrom,
	})
	if err != nil {
		return Account{}, Account{}, err
	}

	accountTo, err := q.UpdateAccountBalance(ctx, UpdateAccountBalanceParams{
		ID:     accountIdTo,
		Amount: amountTo,
	})
	if err != nil {
		return Account{}, Account{}, err
	}

	return accountFrom, accountTo, nil
}
