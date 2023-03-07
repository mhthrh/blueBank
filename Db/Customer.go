package Db

import (
	"context"
	"fmt"
	"github.com/mhthrh/BlueBank/Entity"
	"github.com/pkg/errors"
)

func (d *dataBase) Create(ctx context.Context, customer Entity.Customer) error {
	cnt, err := d.Exist(ctx, customer.UserName)
	if err != nil {
		return fmt.Errorf("customer name problem, %w", err)
	}
	if cnt != 0 {
		return fmt.Errorf("userExist, %w", err)
	}
	c.Text = customer.PassWord
	result, err := d.db.ExecContext(ctx, fmt.Sprintf("INSERT INTO public.customers(full_name, user_name, hash_password, email)VALUES ('%s', '%s', '%s', '%s')", customer.FullName, customer.UserName, c.Sha256(), customer.Email))
	if err != nil {
		return errors.Wrap(err, "cannot insert to db.customer")
	}
	count, _ := result.RowsAffected()
	if count != 1 {
		return errors.Wrap(err, "insert issue to db.customer")
	}
	return nil
}

func (d *dataBase) Login(ctx context.Context, login Entity.CustomerLogin) (*Entity.Customer, error) {
	var customer Entity.Customer
	c.Text = login.PassWord
	row := d.db.QueryRowContext(ctx, fmt.Sprintf("select id, full_name, user_name, email, expires_at, created_at FROM public.customers where user_name='%s' and hash_password='%s'", login.UserName, c.Sha256()))

	err := row.Scan(&customer.ID, &customer.FullName, &customer.UserName, &customer.Email, &customer.CreateAt, &customer.ExpireAt)
	if err != nil {
		return nil, fmt.Errorf("userName/password incorect, %w", err)
	}

	_, err = d.db.ExecContext(ctx, fmt.Sprintf("INSERT INTO public.customer_log( user_name, status)VALUES ('%s', true)", login.UserName))
	if err != nil {
		return nil, fmt.Errorf("db.user_log insert error, %w", err)
	}
	return &customer, nil
}

func (d *dataBase) Exist(ctx context.Context, userName string) (int, error) {
	var count int
	rows, err := d.db.QueryContext(ctx, fmt.Sprintf("select count(*) FROM public.customers where user_name='%s'", userName))
	if err != nil {
		return 0, fmt.Errorf("db error, %w", err)
	}
	for rows.Next() {
		err = rows.Scan(&count)
		if err != nil {
			return 0, fmt.Errorf("db fetch error, %w", err)
		}
	}

	return count, nil
}
