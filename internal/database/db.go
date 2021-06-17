package database

import "github.com/jackc/pgx"

type DB *pgx.Conn

func GetDB(host string, port int64, database, user, password string) (DB, error) {
	c, err := pgx.Connect(pgx.ConnConfig{
		Host:                 host,
		Port:                 uint16(port),
		Database:             database,
		User:                 user,
		Password:             password,
		//Logger:               nil, //TODO db logger
		//LogLevel:             pgx.LogLevelError,
	})
	if err != nil {
		return c, err
	}

	return c, nil
}