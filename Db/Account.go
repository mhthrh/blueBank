package Db

import (
	"context"
	"fmt"
	"github.com/mhthrh/BlueBank/Entity"
	"github.com/pkg/errors"
)

func (d *dataBase) CreateAccount(ctx context.Context, account Entity.Account) error {
	count, err := d.ExistAccount(ctx, account)
	if err != nil {
		return fmt.Errorf("cannot create account, %w", err)
	}
	if count != 0 {
		return fmt.Errorf("douplicate account")
	}

	result, err := d.db.ExecContext(ctx, fmt.Sprintf("INSERT INTO public.accounts(customer_id, account_no, balance)VALUES ('%d', '%d', '%d')", account.CustomerId, account.AccountNumber, account.Balance))
	if err != nil {
		return errors.Wrap(err, "cannot insert to db.accounts, ")
	}
	cnt, _ := result.RowsAffected()
	if cnt != 1 {
		return errors.Wrap(err, "insert issue to db.accounts")
	}

	return nil
}
func (d *dataBase) BalanceAccount(ctx context.Context, account Entity.Account) (int64, error) {
	var balance int64
	rows, err := d.db.QueryContext(ctx, fmt.Sprintf("select balance FROM public.customers where account_no='%d'", account.AccountNumber))
	if err != nil {
		return 0, fmt.Errorf("db error, %w", err)
	}
	err = rows.Scan(&balance)
	if err != nil {
		return 0, fmt.Errorf("db fetch error, %w", err)
	}

	return balance, nil
}
func (d *dataBase) LockAmount(ctx context.Context, account Entity.Account) (e error) {
	commit := false
	var balance int64
	tx, err := d.db.Begin()
	defer func() {
		if commit {
			e = tx.Commit()
		} else {
			e = tx.Rollback()
		}
	}()
	row := tx.QueryRowContext(ctx, fmt.Sprintf("select balance-balance_lock from public.accounts where account_no='%d' FOR UPDATE", account.AccountNumber))
	if err = row.Scan(&balance); err != nil {
		return fmt.Errorf("canot lock some amount, %w", err)
	}
	res, err := tx.ExecContext(ctx, fmt.Sprintf(" UPDATE public.accounts SET  balance=balance-'%d', balance_lock=balance_lock+'%d' WHERE account_no='%d'", account.LockAmount, account.LockAmount, account.AccountNumber))
	if err != nil {
		return fmt.Errorf("canot update lock amount, %w", err)
	}
	cnt, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("canot get result lock amount, %w", err)
	}
	if cnt != 1 {
		return fmt.Errorf("canot update lock amount, %w", err)
	}
	return nil
}
func (d *dataBase) ExistAccount(ctx context.Context, account Entity.Account) (int, error) {

	var count int
	row := d.db.QueryRowContext(ctx, fmt.Sprintf("select count(*) FROM public.accounts where account_no='%d'", account.AccountNumber))

	err := row.Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("db error, %w", err)
	}

	return count, nil
}
