package api

import (
	"context"

	"github.com/jackc/pgx"
	"github.com/jackc/pgx/v5/pgconn"
)

type PosgresInterface interface {
	Begin(context.Context) (pgx.Tx, error)
	Exec(context.Context, string, ...interface{}) (pgconn.CommandTag, error)
	QueryRow(context.Context, string, ...interface{}) pgx.Row
	Query(context.Context, string, ...interface{}) (pgx.Rows, error)
	Ping(context.Context) error
	Close()
}

type Database struct {
    DB PosgresInterface
}

func NewDatabase(ds PosgresInterface) Database {
    return Database{DB: ds}
}