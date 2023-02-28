package Db

import (
	"context"
	"fmt"
	"github.com/mhthrh/BlueBank/Entity"
)

func (d *dataBase) NewBilan(ctx context.Context, b Entity.Bilan) error {
	if _, err := d.ExistBilan(ctx, b); err != nil {
		return fmt.Errorf("conot create bilan, %w", err)
	}
	if (b.IsCredit && b.Amount <= 0) || (!b.IsCredit && b.Amount > 0) {
		return fmt.Errorf("type mismatch")
	}
	_, err := d.db.ExecContext(ctx, fmt.Sprintf("INSERT INTO public.bilans(bilan, balance, is_credit) VALUES ('%d', '%d', '%t')", b.Number, b.Amount, b.IsCredit))
	return err
}

func (d *dataBase) BalanceBilan(ctx context.Context, b Entity.Bilan) (int64, error) {
	var amount int64
	rows, err := d.db.QueryContext(ctx, fmt.Sprintf("select balance FROM public.bilans where bilan='%d'", b.Number))
	if err != nil {
		return 0, fmt.Errorf("db error, %w", err)
	}
	for rows.Next() {
		err = rows.Scan(&amount)
		if err != nil {
			return amount, fmt.Errorf("db fetch error, %w", err)
		}
	}

	return amount, nil
}
func (d *dataBase) ExistBilan(ctx context.Context, b Entity.Bilan) (bool, error) {
	var count int
	rows, err := d.db.QueryContext(ctx, fmt.Sprintf("select count(*) FROM public.bilans where bilan='%d'", b.Number))
	if err != nil {
		return false, fmt.Errorf("db error, %w", err)
	}
	for rows.Next() {
		err = rows.Scan(&count)
		if err != nil {
			return false, fmt.Errorf("db fetch error, %w", err)
		}
	}

	return true, nil
}
