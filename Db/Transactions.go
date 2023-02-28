package Db

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/mhthrh/BlueBank/Entity"
)

func (d *dataBase) Deposit(ctx context.Context, transactions ...Entity.Transaction) error {
	commit := false
	tx, err := d.db.Begin()
	if err != nil {
		return fmt.Errorf("cannot create transaction")
	}
	defer func() {
		if commit {
			_ = tx.Commit()
			return
		}
		_ = tx.Rollback()
	}()
	var id uuid.UUID
	for _, trans := range transactions {
		if trans.Credit.Amount != trans.Debit.Amount {
			return fmt.Errorf("amount mismatch")
		}
		id, _ = uuid.NewRandom()
		// debit account
		result, err := tx.ExecContext(ctx, fmt.Sprintf("update public.accounts SET  balance=balance-'%d' where account_no='%d'", trans.Debit.Amount, trans.Debit.AccountNumber))
		if err != nil {
			return fmt.Errorf("canot deposit account, %w", err)
		}
		cnt, err := result.RowsAffected()
		if err != nil || cnt != 1 {
			return fmt.Errorf("canot deposit account, %w", err)
		}
		result, err = tx.ExecContext(ctx, fmt.Sprintf("INSERT INTO public.transactions(id, bilan_id, account_id, amount)VALUES ('%s','%d','%d','%d')", id.String(), trans.Credit.Bilan, trans.Debit.AccountNumber, -trans.Debit.Amount))
		if err != nil {
			return fmt.Errorf("canot deposit account, %w", err)
		}
		//credit bilan
		result, err = tx.ExecContext(ctx, fmt.Sprintf("update public.bilans SET  balance=balance+'%d' where bilan='%d'", trans.Credit.Bilan, trans.Credit.AccountNumber))
		if err != nil {
			return fmt.Errorf("canot deposit account, %w", err)
		}
		cnt, err = result.RowsAffected()
		if err != nil || cnt != 1 {
			return fmt.Errorf("canot deposit account, %w", err)
		}
		result, err = tx.ExecContext(ctx, fmt.Sprintf("INSERT INTO public.transactions(id, bilan_id, account_id, amount)VALUES ('%s','%d','%d','%d')", id.String(), trans.Credit.Bilan, trans.Debit.AccountNumber, trans.Debit.Amount))
		if err != nil {
			return fmt.Errorf("canot deposit account, %w", err)
		}

	}
	commit = true
	return nil
}

func (d *dataBase) Withdraw(ctx context.Context, transactions ...Entity.Transaction) error {
	commit := false
	tx, err := d.db.Begin()
	if err != nil {
		return fmt.Errorf("cannot create transaction")
	}
	defer func() {
		if commit {
			_ = tx.Commit()
			return
		}
		_ = tx.Rollback()
	}()
	var id uuid.UUID
	for _, trans := range transactions {
		if trans.Credit.Amount != trans.Debit.Amount {
			return fmt.Errorf("amount mismatch")
		}
		id, _ = uuid.NewRandom()
		// credit account
		result, err := tx.ExecContext(ctx, fmt.Sprintf("update public.accounts SET  balance=balance+'%d' where account_no='%d'", trans.Debit.Amount, trans.Debit.AccountNumber))
		if err != nil {
			return fmt.Errorf("canot deposit account, %w", err)
		}
		cnt, err := result.RowsAffected()
		if err != nil || cnt != 1 {
			return fmt.Errorf("canot deposit account, %w", err)
		}
		result, err = tx.ExecContext(ctx, fmt.Sprintf("INSERT INTO public.transactions(id, bilan_id, account_id, amount)VALUES ('%s','%d','%d','%d')", id.String(), trans.Credit.Bilan, trans.Debit.AccountNumber, trans.Debit.Amount))
		if err != nil {
			return fmt.Errorf("canot deposit account, %w", err)
		}
		//debit bilan
		result, err = tx.ExecContext(ctx, fmt.Sprintf("update public.bilans SET  balance=balance-'%d' where bilan='%d'", trans.Credit.Bilan, trans.Credit.AccountNumber))
		if err != nil {
			return fmt.Errorf("canot deposit account, %w", err)
		}
		cnt, err = result.RowsAffected()
		if err != nil || cnt != 1 {
			return fmt.Errorf("canot deposit account, %w", err)
		}
		result, err = tx.ExecContext(ctx, fmt.Sprintf("INSERT INTO public.transactions(id, bilan_id, account_id, amount)VALUES ('%s','%d','%d','%d')", id.String(), trans.Credit.Bilan, trans.Debit.AccountNumber, -trans.Debit.Amount))
		if err != nil {
			return fmt.Errorf("canot deposit account, %w", err)
		}

	}
	commit = true
	return nil
}
