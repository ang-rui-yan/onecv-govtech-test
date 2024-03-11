package api

import (
	"context"
	"database/sql"
	"errors"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)


type DBStore struct {
	dbPool *pgxpool.Pool
}
func NewDBStore(pool *pgxpool.Pool) DBStore {
	return DBStore{dbPool: pool}
}

func (ds DBStore) Pool() *pgxpool.Pool {
	return ds.dbPool
}

func NewDBPool(dbURL string) (*pgxpool.Pool, func(), error) {

	f := func() {}

    // create pgx connection pool
	pool, err := pgxpool.New(context.Background(), dbURL)

	if err != nil {
        return nil, f, errors.New("database connection error")
	}

	err = validateDBPool(pool)

	if err != nil {
		return nil, f, err
	}

    // return connection pool and inline function to close/ clear the pool if not used. 
    // return nil for the error since there should be no error to this point
	return pool, func() { pool.Close() }, nil
}

// validates if the pool is open
func validateDBPool(pool *pgxpool.Pool) error {
    err := pool.Ping(context.Background())

	if err != nil {
        return errors.New("database connection error")
	}

	var (
		currentDatabase string
		currentUser     string
		dbVersion       string
	)
	
    sqlStatement := `select current_database(), current_user, version();`
	row := pool.QueryRow(context.Background(), sqlStatement)
	err = row.Scan(&currentDatabase, &currentUser, &dbVersion)

	switch {
		case err == sql.ErrNoRows:
			return errors.New("no rows were returned")
		case err != nil:
			return errors.New("database connection error")
		default:
			log.Printf("database version: %s\n", dbVersion)
			log.Printf("current database user: %s\n", currentUser)
			log.Printf("current database: %s\n", currentDatabase)
	}

	return nil
}