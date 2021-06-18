package database

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4"
)

func GetDB(ctx context.Context, host string, port int64, database, user, password string) (*pgx.Conn, error) {
	c, err := pgx.Connect(ctx, fmt.Sprintf("postgres://%v:%v@%v:%v/%v", user, password, host, port, database))
	if err != nil {
		return c, err
	}

	return c, nil
}