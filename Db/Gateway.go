package Db

import (
	"context"
	"fmt"
	"github.com/mhthrh/BlueBank/Entity"
)

func (d *dataBase) GatewayLogin(ctx context.Context, gateway Entity.GatewayLogin) (*Entity.Gateway, error) {
	var g Entity.Gateway
	c.Text = gateway.Password
	rows, err := d.db.QueryContext(ctx, fmt.Sprintf("SELECT gateway_name,user_name, ips, status FROM public.gateways where status=true and user_name='%s' and hash_password='%s'", gateway.UserName, c.Sha256()))
	if err != nil {
		return nil, fmt.Errorf("db error, %w", err)
	}
	s := false
	for rows.Next() {
		err = rows.Scan(&g.GatewayName, &g.UserName, &g.Ips, &g.Status)
		if err != nil {
			return nil, fmt.Errorf("db fetch error, %w", err)
		}
		s = true
	}
	if !s {
		return nil, fmt.Errorf("gateway not found")
	}
	return &g, nil
}
