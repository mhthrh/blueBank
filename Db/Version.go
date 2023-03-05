package Db

import (
	"context"
	"fmt"
)

func (d *dataBase) GetVersion(ctx context.Context, key string) (string, error) {
	var version string
	row := d.db.QueryRowContext(ctx, fmt.Sprintf("select value from public.config where key='%s'", key))
	err := row.Scan(&version)
	return version, err
}
